package goPsdLib

import (
	"image/png"
	"log"
	"os"
)

func ParseFormFile(filename string) {
	document := NewDocument(filename)

	_fileHeader := parseFileHeader(document)
	parseColorModeData(document)
	parseImageResources(document)
	parseLayerMask(document, _fileHeader)
	parseImageData(document, _fileHeader)
}

func parseFileHeader(doc *document) *FileHeader {
	file_header := &FileHeader{}
	file_header.Inits(doc)
	log.Printf("parseFileHeader result :%+v \n ", file_header)
	return file_header
}

func parseColorModeData(doc *document) *ColorModeData {
	color_mode_data := &ColorModeData{}
	color_mode_data.Inits(doc)
	// log.Printf("parseColorModeData result :%+v \n ", color_mode_data)
	return color_mode_data
}

func parseImageResources(doc *document) *ImageResources {
	image_resources := &ImageResources{}
	image_resources.Inits(doc)

	// log.Printf("image_resources result :%+v \n ", image_resources)
	return image_resources
}

func parseLayerMask(doc *document, _fileHeader *FileHeader) *LayerMaskInformation {
	layerMask := &LayerMaskInformation{}
	layerMask.Inits(doc, _fileHeader)

	layerMask.LayerInfo.GetImage()

	return layerMask
}

func parseImageData(doc *document, _fileHeader *FileHeader) *ImageData {
	img := &ImageData{}
	img.Inits(doc, _fileHeader)

	f, _ := os.Create("./images/psd.png")

	png.Encode(f, img.Image)

	return img
}
