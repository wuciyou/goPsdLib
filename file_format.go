package goPsdLib

import (
	"image"
	"image/color"
	"log"
)

func (header *FileHeader) Inits(document *document) {
	header.FileType = FileType(document.Next(4))
	document.readUint16(&header.Version)

	header.Reserved = document.Next(6)

	document.readUint16(&header.Channels)
	document.readUint32(&header.Height)
	document.readUint32(&header.Width)
	document.readUint16(&header.Depth)

	var color_model_uint16 uint16
	document.readUint16(&color_model_uint16)

	header.ColorMode = colorMode(color_model_uint16)
}

func (colorMode *ColorModeData) Inits(doc *document) {
	doc.readUint32(&colorMode.Length)
	colorMode.ColorData = doc.Next(int(colorMode.Length))
}

func (imageResources *ImageResources) Inits(doc *document) {
	doc.readUint32(&imageResources.Length)

	imageResources.ImageResourcesBuf = NewDocumentFromByte(doc.Next(int(imageResources.Length)))

	imageResources.ImageResource = make(map[uint16]interface{})

	for {

		signature := string(imageResources.ImageResourcesBuf.Next(4))

		if signature != "8BIM" {
			log.Panicf("Wrong signature of resource %s \n", signature)
		}

		var identifier uint16
		imageResources.ImageResourcesBuf.readUint16(&identifier)

		var name string
		imageResources.ImageResourcesBuf.readPascalString(&name)

		var actualDataSize uint32
		imageResources.ImageResourcesBuf.readUint32(&actualDataSize)

		resource_data_buf := NewDocumentFromByte(imageResources.ImageResourcesBuf.Next(int(actualDataSize)))

		// log.Printf("identifier:%d name:%s , actualDataSize:%d , resource_data:%v \n ", identifier, name, actualDataSize, resource_data_buf.Bytes())

		switch identifier {
		case 1050:
			imageResources.ImageResource[identifier] = imageResourceIDs1050(resource_data_buf, identifier)
		case 1057:
			version_info := &ResourceVersionInfo{}

			resource_data_buf.readUint32(&version_info.Version)

			version_info.HasRealMergedData, _ = resource_data_buf.ReadByte()

			resource_data_buf.readUnicodeString(&version_info.WriterName)
			resource_data_buf.readUnicodeString(&version_info.ReaderName)
			resource_data_buf.readUint32(&version_info.FileVersion)

			// log.Printf("id:%d, version_info:%+v  \n  ", identifier, version_info)
			imageResources.ImageResource[identifier] = version_info
		default:
			log.Printf("Unknown identifier:%d", identifier)
		}

		if actualDataSize%2 != 0 {
			imageResources.ImageResourcesBuf.Next(1)
		}

		if imageResources.ImageResourcesBuf.Len() <= 0 {
			break
		}
	}
}

func (layreMask *LayerMaskInformation) Inits(doc *document, _fileHeader *FileHeader) {

	if _fileHeader.FileType == FILETYPE_PSB {
		doc.readUint64(&layreMask.Length)
	} else {
		var lengthU32 uint32
		doc.readUint32(&lengthU32)
		layreMask.Length = uint64(lengthU32)
	}

	layerMaskInformationBuf := NewDocumentFromByte(doc.Next(int(layreMask.Length)))

	layer := &LayerInfo{}

	layerMaskInformationBuf.isPSB = _fileHeader.FileType == FILETYPE_PSB

	layer.readStructure(layerMaskInformationBuf)

	layreMask.LayerInfo = layer
}

// #http://www.adobe.com/devnet-apps/photoshop/fileformatashtml/#50577409_89817
func (img *ImageData) Inits(doc *document, _fileHeader *FileHeader) {
	doc.readUint16(&img.CompressionMethod)

	width := int(_fileHeader.Width)
	height := int(_fileHeader.Height)
	channels := int(_fileHeader.Channels)
	byteCounts := make([]int16, channels*height)
	isRLE := img.CompressionMethod == 1

	if isRLE {
		for i := range byteCounts {
			var dataU16 uint16
			doc.readUint16(&dataU16)
			byteCounts[i] = int16(dataU16)
		}
	}

	chanData := make(map[int][]byte)
	for i := 0; i < channels; i++ {
		var data []byte
		if isRLE {
			chanIndex := i * height
			for j := 0; j < height; j++ {
				length := byteCounts[chanIndex]
				chanIndex++
				dataTemp := doc.Next(int(length))
				// dataTemp := doc.ReadSignedBytes(int(length))
				data = append(data, UnpackRLEBits(dataTemp, width)...)
			}
		} else {
			data = doc.Next(width * height)
		}
		chanData[i] = data
	}

	img.Image = image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			i := x + (y * width)
			red := byte(chanData[0][i])
			green := byte(chanData[1][i])
			blue := byte(chanData[2][i])
			img.Image.Set(x, y, color.RGBA{red, green, blue, 255})
		}
	}
}
