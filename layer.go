package goPsdLib

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func (layer *LayerInfo) GetImage() {
	for key_index, layerRecords := range layer.LayerRecords {
		width := int(layerRecords.Rectangle.Width)
		height := int(layerRecords.Rectangle.Height)
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		switch layerRecords.ChannelsNumber {
		case 3:
		case 4, 5:
			c := layerRecords.ChannelInformations
			for x := 0; x < width; x++ {
				for y := 0; y < height; y++ {
					i := x + (y * width)
					red := byte(c[1].Data[i])
					green := byte(c[2].Data[i])
					blue := byte(c[3].Data[i])
					alpha := byte(c[0].Data[i])
					img.Set(x, y, color.RGBA{red, green, blue, alpha})
				}
			}
		}

		img_file_name := fmt.Sprintf("./images/_%d_%s.png", key_index, layerRecords.LayerName)
		f_png, err := os.Create(img_file_name)
		if err != nil {
			log.Printf("Can't create file[%s], err:%v \n ", img_file_name, err)
		} else {
			png.Encode(f_png, img)
		}

	}
}
