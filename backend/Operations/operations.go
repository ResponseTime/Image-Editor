package Operations

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"main/Db"
	"main/Models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	mu         sync.Mutex
	userStacks = make(map[string]*Models.UserStack)
	Client     = Db.GetClient()
)

func moveFile(sourcePath, destinationPath string) error {
	err := os.Rename(sourcePath, destinationPath)
	if err != nil {
		return err
	}
	return nil
}
func Upload(c *gin.Context) {
	mu.Lock()
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
	}
	_, exist := userStacks[email.(string)]
	if exist {
		ImagePath := userStacks[email.(string)].CurrentImage.Path
		err := os.Remove(ImagePath)
		if err != nil {
			fmt.Printf("Error removing %s", ImagePath)
		}
	}
	userStacks[email.(string)] = &Models.UserStack{}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Dir, _ := os.Getwd()
	err = c.SaveUploadedFile(file, Dir+file.Filename)
	userStacks[email.(string)].CurrentImage = &Models.Image{Path: Dir + file.Filename}
	userStacks[email.(string)].UndoStack = []*Models.Image{}
	userStacks[email.(string)].RedoStack = []*Models.Image{}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
	defer mu.Unlock()
}

func Rotate(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
	}
	angle, err := strconv.ParseFloat(c.PostForm("angle"), 64)
	if err != nil {
		log.Fatal(err)
	}
	_, exist := userStacks[email.(string)]
	if exist {
		imagePath := userStacks[email.(string)].CurrentImage.Path
		file, err := os.Open(imagePath)
		if err != nil {
			log.Fatal(err)
		}
		img, err := imaging.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		rotatedImg := imaging.Rotate(img, angle, color.Transparent)
		err = imaging.Save(rotatedImg, imagePath)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"Success": "True"})
		file.Close()
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}

func Save(c *gin.Context) {
	email, exists := c.Get("email")
	pn := c.Param("pname")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
	}
	_, exist := userStacks[email.(string)]
	if exist {
		imageFile := userStacks[email.(string)].CurrentImage.Path
		data := bson.M{
			"ImagePath":   imageFile,
			"User":        email.(string),
			"ProjectName": pn,
			"Date":        fmt.Sprintf("Current Date: %04d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day()),
			"TimeCreated": time.Now().UTC().Format(time.TimeOnly),
		}
		Dir, _ := os.Getwd()
		moveFile(imageFile, Dir+"/uploads"+pn)
		insertRes, err := Client.Database("users").Collection("image_details").InsertOne(context.Background(), data)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"success": insertRes.InsertedID})

	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}

func Export(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
	}
	_, exist := userStacks[email.(string)]
	if exist {
		imageFile := userStacks[email.(string)].CurrentImage.Path
		filesplit := strings.Split(userStacks[email.(string)].CurrentImage.Path, "/")
		extsplit := strings.Split(userStacks[email.(string)].CurrentImage.Path, ".")
		filename, extType := filesplit[len(filesplit)-1], extsplit[1]
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Writer.Header().Set("filename", filename)
		c.Writer.Header().Set("Content-Type", fmt.Sprintf("image/%s", extType))
		c.File(imageFile)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}
func GetDetails(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
	}
	filter := bson.M{"User": email}
	var res []Models.Result
	var result Models.Result
	cursor, err := Client.Database("users").Collection("image_details").Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		res = append(res, result)

	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}
