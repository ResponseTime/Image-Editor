package Auth

import (
	"context"
	"fmt"
	"log"
	"main/Db"
	"main/Models"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func CreateToken(user *Models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func AuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	email, ok := claims["email"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from token claims"})
		c.Abort()
		return
	}
	c.Set("email", email)
	c.Next()
}

func Signup(c *gin.Context) {
	var user Models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if Db.Client == nil {
		log.Fatal("Connect To Database")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
	if err != nil {
		log.Fatal(err)
	}
	insertRes, err := Db.Client.Database("users").Collection("login_details").InsertOne(context.Background(), bson.M{"email": user.Email, "password": hashedPassword})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted document with ID %v\n", insertRes.InsertedID)
	c.JSON(http.StatusOK, gin.H{"success": insertRes.InsertedID})

}

func Login(c *gin.Context) {
	var user Models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if Db.Client == nil {
		log.Fatal("Connect To Database")
	}
	var result Models.User
	filter := bson.M{"email": user.Email}
	err := Db.Client.Database("users").Collection("login_details").FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
	if err != nil {
		log.Fatal(err)
	}
	ok := bcrypt.CompareHashAndPassword(hashedPassword, []byte(user.Password))
	if ok != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password dont match"})
	} else {
		tokenString, err := CreateToken(&user)
		if err != nil {
			fmt.Println("Error creating token:", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}
