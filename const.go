package goPsdLib

type colorMode uint16

const (
	BIT_MAP       colorMode = 0
	GRAYSCALE               = 1
	INDEXED                 = 2
	RGB                     = 3
	CMYK                    = 4
	MULTI_CHANNEL           = 7
	DUOTONE                 = 8
	LAB                     = 9
)

type FileType string

const (
	FILETYPE_BBPS FileType = "8BPS"
	FILETYPE_PSB           = "PSB"
)
