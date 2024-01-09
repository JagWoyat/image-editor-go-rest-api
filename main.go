package main

import (
	"image-editor/models"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	server.Use(cors.New(config))

	server.GET("/image/:path", getImage)
	server.POST("/image", uploadImage)

	server.Run(":4000")
}

func getImage(context *gin.Context) {
	filename := context.Param("path")
	f, err := os.Open("output/" + filename)
	if err != nil {
		return
	}
	defer f.Close()

	context.File("output/" + filename)
}

func uploadImage(context *gin.Context) {

	var req models.Request

	if err := context.ShouldBind(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// file, err := context.FormFile("img")

	// if err != nil {
	// 	context.JSON(http.StatusBadRequest, gin.H{"message": "Could not upload image"})
	// 	return
	// }

	var uploadedImage models.Image

	// name := context.Request.FormValue("name")

	// extension := filepath.Ext(file.Filename)
	uploadedImage.SetName(req.Image.Filename)

	err := context.SaveUploadedFile(req.Image, "images/"+uploadedImage.Name)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": uploadedImage.Name})
		return
	}

	err = uploadedImage.OpenImage("images/" + uploadedImage.Name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode image"})
		return
	}

	if req.Height+req.Width > 0 {
		uploadedImage.ScaleImage(req.Width, req.Height)
	}

	if req.MirroredHor {
		uploadedImage.MirrorImageHorizontaly()
	}

	if req.MirroredVer {
		uploadedImage.MirrorImageVerticaly()
	}

	if req.Grayscale {
		uploadedImage.GrayscaleImage()
	}

	uploadedImage.Path = "output/" + uploadedImage.Name

	uploadedImage.SaveImage()

	context.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully", "name": uploadedImage.Name})

}
