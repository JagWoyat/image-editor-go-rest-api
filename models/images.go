package models

import (
	"image"
	"mime/multipart"

	"github.com/disintegration/imaging"
)

type Request struct {
	Name        string                `form:"name" binding:"required"`
	Image       *multipart.FileHeader `form:"image" binding:"required"`
	MirroredHor bool                  `form:"mirrored_h"`
	MirroredVer bool                  `form:"mirrored_v"`
	Grayscale   bool                  `form:"grayscale"`
	Height      int                   `form:"height"`
	Width       int                   `form:"width"`
}

type Image struct {
	Name  string
	Path  string
	Image image.Image
}

func (img *Image) SetName(newName string) {
	img.Name = newName
}

func (img *Image) SetPath(newPath string) {
	img.Path = newPath
}

func (img *Image) OpenImage(path string) (err error) {
	var image image.Image
	image, err = imaging.Open(path)
	if err != nil {
		return err
	}
	img.Image = image

	return nil
}

func (img *Image) SaveImage() (err error) {
	err = imaging.Save(img.Image, img.Path)
	if err != nil {
		return err
	}
	return nil
}

func (img *Image) ScaleImage(width, height int) {
	img.Image = imaging.Resize(img.Image, width, height, imaging.Lanczos)
}

func (img *Image) MirrorImageHorizontaly() {
	img.Image = imaging.FlipH(img.Image)
}

func (img *Image) MirrorImageVerticaly() {
	img.Image = imaging.FlipV(img.Image)
}

func (img *Image) GrayscaleImage() {
	img.Image = imaging.Grayscale(img.Image)
	img.Image = imaging.AdjustContrast(img.Image, 20)
	img.Image = imaging.Sharpen(img.Image, 2)
}
