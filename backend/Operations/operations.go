package Operations

import (
	"context"
	"fmt"
	"io"
	"log"
	"main/Db"
	"main/Models"
	"main/Utilities"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"gocv.io/x/gocv"
)

var (
	mu         sync.Mutex
	userStacks = make(map[string]*Models.UserStack)
)

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
	err = c.SaveUploadedFile(file, Dir+"/tempupload/"+file.Filename)
	userStacks[email.(string)].CurrentImage = &Models.Image{Path: Dir + "/tempupload/" + file.Filename}
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
	_, exist := userStacks[email.(string)]
	if exist {
		img := gocv.IMRead(userStacks[email.(string)].CurrentImage.Path, gocv.IMReadColor)
		if img.Empty() {
			fmt.Printf("Failed to read image: %s\n", userStacks[email.(string)].CurrentImage.Path)
			os.Exit(1)
		}
		rotated := gocv.NewMat()
		gocv.Rotate(img, &rotated, gocv.Rotate90Clockwise)
		outPath := filepath.Join(userStacks[email.(string)].CurrentImage.Path)
		if ok := gocv.IMWrite(outPath, rotated); !ok {
			fmt.Printf("Failed to write image: %s\n")
			os.Exit(1)
		}
		c.JSON(http.StatusOK, gin.H{"Success": "True"})
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
			"Date":        fmt.Sprintf("%04d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day()),
			"TimeCreated": time.Now().UTC().Format(time.TimeOnly),
		}
		Dir, _ := os.Getwd()
		Utilities.MoveFile(imageFile, Dir+"\\uploads\\"+pn)
		insertRes, err := Db.Client.Database("users").Collection("image_details").InsertOne(context.Background(), data)
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
	cursor, err := Db.Client.Database("users").Collection("image_details").Find(context.Background(), filter)
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

func Bright_inc(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
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
		brightImg := imaging.AdjustBrightness(img, 10)
		err = imaging.Save(brightImg, imagePath)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"Success": "True"})
		file.Close()
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}

func Bright_dec(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
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
		brightImg := imaging.AdjustBrightness(img, -10)
		err = imaging.Save(brightImg, imagePath)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"Success": "True"})
		file.Close()
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}

func Contrast_inc(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
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
		contrastImg := imaging.AdjustContrast(img, 15)
		err = imaging.Save(contrastImg, imagePath)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"Success": "True"})
		file.Close()
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}

func Contrast_dec(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
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
		contrastImg := imaging.AdjustContrast(img, -15)
		err = imaging.Save(contrastImg, imagePath)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"Success": "True"})
		file.Close()
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}

func Grayscale(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
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
		grayImg := imaging.Grayscale(img)
		err = imaging.Save(grayImg, imagePath)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"Success": "True"})
		file.Close()
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}

func Resize(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
	}
	var dimensions Models.Resizer
	if err := c.ShouldBindJSON(&dimensions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
		resImg := imaging.Resize(img, dimensions.Height, dimensions.Width, imaging.Lanczos)
		err = imaging.Save(resImg, imagePath)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"Success": "True"})
		file.Close()
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}

func GetImage(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
	}
	_, exist := userStacks[email.(string)]
	if exist {
		imageFile, err := os.Open(userStacks[email.(string)].CurrentImage.Path)
		extType := strings.Split(userStacks[email.(string)].CurrentImage.Path, ".")[1]
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening image file"})
			return
		}
		c.Header("Content-Type", fmt.Sprintf("image/%s", extType))
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", userStacks[email.(string)].CurrentImage.Path))
		_, err = io.Copy(c.Writer, imageFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error copying image to response"})
			return
		}
		imageFile.Close()
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}

func Undo(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
	}
	if exists {
		if len(userStacks[email.(string)].UndoStack) == 10 {
			userStacks[email.(string)].UndoStack = userStacks[email.(string)].UndoStack[1:]
		}
		userStacks[email.(string)].UndoStack = append(userStacks[email.(string)].UndoStack, userStacks[email.(string)].CurrentImage)
		c.JSON(http.StatusOK, gin.H{"success": true})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}

func Redo(c *gin.Context) {
	_, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return
	}
	if exists {

	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload a image first"})
	}
}
