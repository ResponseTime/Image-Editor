package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"image/color"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/disintegration/imaging"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type User struct {
	Email string `json:"email"`
	Password  string `json:"password"`
}
type userStack struct {
	UndoStack    []*Image
	RedoStack    []*Image
	CurrentImage *Image
}
type Image struct {
	Path string
}
var (
	userStacks = make(map[string]*userStack)
	mu sync.Mutex
	client1 *mongo.Client
)

func connectDb(){
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client,err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	client1 = client
}

func createToken(user *User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix()
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func authMiddleware(c *gin.Context) {
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

func signup(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	insertRes,err := client1.Database("users").Collection("login_details").InsertOne(context.Background(),user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Inserted document with ID %v\n", insertRes.InsertedID)
	c.JSON(http.StatusOK, gin.H{"success": insertRes.InsertedID})
}
func login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var result User
	filter := bson.M{"email": user.Email}
	err := client1.Database("users").Collection("login_details").FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user.Email)
	fmt.Println(result.Email)
	if result.Password == user.Password {
		tokenString,err:= createToken(&user);
		if err != nil {
			fmt.Println("Error creating token:", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}else{
		c.JSON(http.StatusBadRequest, gin.H{"error": "password dont match"})
	}
}

func upload(c *gin.Context) {
	mu.Lock()
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
	}
	_,exist:= userStacks[email.(string)];
	if !exist {
		userStacks[email.(string)] = &userStack{}
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	 	}
		err = c.SaveUploadedFile(file, "tempupload/"+file.Filename)
		userStacks[email.(string)].CurrentImage = &Image{Path: "tempupload/"+file.Filename}
		userStacks[email.(string)].UndoStack = []*Image{}
		userStacks[email.(string)].RedoStack = []*Image{}
     	if err != nil {
        	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        	return
    	}
		c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
	}else{
		ImagePath := userStacks[email.(string)].CurrentImage.Path
		err:=os.Remove(ImagePath)
		if err != nil{
			fmt.Printf("Error removing %s",ImagePath)
		}else{
			file, err := c.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
	 			}
			err = c.SaveUploadedFile(file, "tempupload/"+file.Filename)
			userStacks[email.(string)].CurrentImage = &Image{Path: "tempupload/"+file.Filename}
			userStacks[email.(string)].UndoStack = []*Image{}
			userStacks[email.(string)].RedoStack = []*Image{}
     		if err != nil {
        		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        		return
    		}
			c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
		}
	}
	defer mu.Unlock()
}


func undo(c *gin.Context) {
	
}
func redo(c *gin.Context) {

}

func crop(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
	}
	_,exist:= userStacks[email.(string)];
	if exist{
		diff:=&Image{Path: "awdawd"}
		userStacks[email.(string)].UndoStack = append(userStacks[email.(string)].UndoStack, diff)
	}else{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}
// func imageCenter(img gocv.Mat) image.Point {
// 	height, width := img.Size()[0], img.Size()[1]
// 	return image.Pt(width/2, height/2)
// }
func rotate(c *gin.Context){
	email,exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
	}
	_,exist:= userStacks[email.(string)];
	if exist{
		imagePath := userStacks[email.(string)].CurrentImage.Path
		file, err := os.Open(imagePath)
		if err != nil {
			log.Fatal(err)
		}
		img, err := imaging.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		rotatedImg := imaging.Rotate(img, 90, color.Transparent)
		err = imaging.Save(rotatedImg, imagePath)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK,gin.H{"Success":"True"})
	}else{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}
func getImage(c *gin.Context){
	email,exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
	}
	_,exist:= userStacks[email.(string)];
	if exist{
		imageFile, err := os.Open(userStacks[email.(string)].CurrentImage.Path)
		extType:= strings.Split(userStacks[email.(string)].CurrentImage.Path, ".")[1]
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening image file"})
			return
		}
		defer imageFile.Close()
		c.Header("Content-Type", fmt.Sprintf("image/%s", extType))
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", userStacks[email.(string)].CurrentImage.Path))
		_, err = io.Copy(c.Writer, imageFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error copying image to response"})
			return
		}
	}else{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}