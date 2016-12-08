package goPsdLib

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"unicode/utf16"
)

type document struct {
	isPSB bool
	*bytes.Buffer
}

func NewDocument(fileName string) *document {
	f, e := os.Open(fileName)
	if e != nil {
		log.Printf("[error] Can't open the file \"%s\", error:%v", fileName, e)
		return nil
	}

	doc := &document{Buffer: &bytes.Buffer{}}
	doc.Buffer.ReadFrom(f)
	return doc
}

func NewDocumentFromByte(b []byte) *document {
	doc := &document{Buffer: bytes.NewBuffer(b)}
	return doc
}

func (doc *document) readRectangle() *Rectangle {
	rectangle := &Rectangle{}
	rectangle.readStructure(doc)
	return rectangle
}

func (doc *document) readBool(s *bool) {
	b, _ := doc.ReadByte()
	*s = b == 1
}
func (doc *document) readUnicodeString(s *string) {
	var n uint32
	doc.readUint32(&n)

	array := make([]uint16, n)
	for i := range array {
		if err := binary.Read(doc, binary.BigEndian, &array[i]); err != nil {
			log.Panicln(err)
		}
	}

	*s = string(utf16.Decode(array))
}

func (doc *document) readPascalString(s *string) {
	var strSizeU8 uint8

	doc.readUint8(&strSizeU8)
	strSize := int(strSizeU8)
	if strSize <= 0 {
		strSize = 1
	}
	*s = string(doc.Next(strSize))
}

func (doc *document) readDynamicString(s *string) {
	var strSizeU32 uint32

	doc.readUint32(&strSizeU32)
	strSize := int(strSizeU32)
	if strSize <= 0 {
		strSize = 4
	}
	*s = string(doc.Next(strSize))
}

func (doc *document) readUint8(u *uint8) {
	data, _ := doc.ReadByte()
	*u = uint8(data)
}

func (doc *document) ReadSignedBytes(lens int) []int8 {
	value := make([]int8, lens)
	binary.Read(doc, binary.BigEndian, value)
	return value
}

func (doc *document) readUint16(u *uint16) {
	data := doc.Next(2)
	var new_data uint16
	for i := 0; i < 2; i++ {
		new_data |= uint16(data[i]) << uint((1-i)*8)
	}
	*u = new_data
}

func (doc *document) readInt16(u *int16) {
	data := doc.Next(2)
	var new_data int16
	for i := 0; i < 2; i++ {
		new_data |= int16(data[i]) << uint((1-i)*8)
	}
	*u = new_data
}

func (doc *document) readUint32(u *uint32) {
	data := doc.Next(4)
	var new_data uint32
	for i := 0; i < 4; i++ {
		new_data |= uint32(data[i]) << uint((3-i)*8)
	}
	*u = new_data
}

func (doc *document) readUint64(u *uint64) {
	data := doc.Next(8)
	var new_data uint64
	for i := 0; i < 8; i++ {
		new_data |= uint64(data[i]) << uint((7-i)*8)
	}
	*u = new_data
}

func (doc *document) readInt64(u *float64) {

	if err := binary.Read(doc, binary.BigEndian, u); err != nil {
		panic(err)
	}

}
