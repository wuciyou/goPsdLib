package main

import (
	"github.com/wuciyou/goPsdLib"
	"log"
)

var filename = "./lcd1602.psd"

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	goPsdLib.ParseFormFile(filename)
}
