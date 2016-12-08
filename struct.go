package goPsdLib

import (
	"image"
)

type Rectangle struct {
	Top                 uint32
	Left                uint32
	Bottom              uint32
	Right               uint32
	X, Y, Width, Height uint32
}
type Color struct {
	Red, Green, Blue, Alpha uint16
}

type FileHeader struct {

	// Mac OS: 8BPS
	// Windows: .PSD
	FileType  FileType
	Version   uint16
	Reserved  []byte
	Channels  uint16
	Height    uint32
	Width     uint32
	Depth     uint16
	ColorMode colorMode
}

type ColorModeData struct {
	Length    uint32
	ColorData []byte
}

type ImageResources struct {
	Length            uint32
	ImageResourcesBuf *document
	ImageResource     map[uint16]interface{}
}

/**
 * 0x0421 1057 (Photoshop 6.0) Version Info. 4 bytes version, 1 byte hasRealMergedData , Unicode string: writer name, Unicode string: reader name, 4 bytes file version.
 */
type ResourceVersionInfo struct {
	Version           uint32
	HasRealMergedData byte
	WriterName        string
	ReaderName        string
	FileVersion       uint32
}

/**
 * 0x041A 1050 (Photoshop 6.0) Slices. See (See Slices resource format)[http://www.adobe.com/devnet-apps/photoshop/fileformatashtml/PhotoshopFileFormats.htm#50577409_19931]
 * Adobe Photoshop 6.0 stores slices information for an image in an image resource block.
 * Adobe Photoshop 7.0 added a descriptor at the end of the block for the individual slice info.
 * Adobe Photoshop CS and later changed to version 7 or 8 and uses a Descriptor to defined the Slices data.
 */
type SlicesResourceFormat struct {
	Header *SlicesHeader
	Blocks []*SlicesResourceBlock
}

type SlicesHeader struct {
	Version      uint32
	Top          uint32
	Leftt        uint32
	Bottom       uint32
	Right        uint32
	GroupName    string
	SlicesNumber uint32
}

type SlicesResourceBlock struct {
	ID                  uint32
	GroupId             uint32
	Origin              uint32
	AssociatedLayerID   uint32
	Name                string
	Type                uint32
	Left                uint32
	Top                 uint32
	Right               uint32
	Bottom              uint32
	URL                 string
	Target              string
	Message             string
	AltTag              string
	CellTextIsHTML      bool
	CellText            string
	HorizontalAlignment uint32
	VerticalAlignment   uint32
	AlphaColor          byte
	Red                 byte
	Green               byte
	Blue                byte
	DescriptorVersion   uint32
	descriptorStructure *DescriptorStructure
}

type DescriptorStructure struct {
	Name        string
	ClassId     string
	ItemsNumber uint32
	Items       map[string]*Descriptor
}

type Descriptor struct {
	isSkipKey bool
	Key       string
	Type      string
	Value     StructureEntity
}

type StructureEntity interface {
	readStructure(doc *document)
}

type ReferenceStructureEntity struct {
	entitys map[string]StructureEntity
}

type PropertyStructureEntity struct {
	Name    string
	ClassID string
	KeyId   string
}
type ClassStructureEntity struct {
	Name    string
	ClassID string
}
type EnumeratedReferenceEntity struct {
	Name    string
	ClassID string
	TypeID  string
	Enum    string
}
type OffsetStructureEntity struct {
	Name    string
	ClassID string
	Value   uint32
}
type LongEntity struct {
	Value uint64
}
type ListStructureEntity struct {
	ItemsNumber    uint32
	DescriptorList []*Descriptor
}

type EnumeratedDescriptorEntity struct {
	Type string
	Enum string
}

type StringStructureEntity struct {
	Value string
}
type BooleanStructureEntity struct {
	isTrue bool
}

/*****************************************************************************
 *
 * LayerMaskInformation
 *
 *****************************************************************************
 */
type LayerMaskInformation struct {
	Length                  uint64
	LayerMaskInformationBuf *document
	LayerInfo               *LayerInfo
}

type LayerInfo struct {
	Length       uint64
	LayerCount   uint16
	LayerRecords []*LayerRecords
}

//http://www.adobe.com/devnet-apps/photoshop/fileformatashtml/#50577409_13084
type LayerRecords struct {
	ID                   uint32
	Rectangle            *Rectangle
	ChannelsNumber       uint16
	ChannelInformations  []*ChannelInformation
	BlendModeSignature   string
	BlendModeKey         string
	Opacity              byte
	Clipping             byte
	Flags                byte
	Filler               byte
	ExtraDataFieldLength uint32

	LayerMaskORAdjustmentLayerData *LayerMaskORAdjustmentLayerData
	LayerBlendingRangesData        *LayerBlendingRangesData
	LayerName                      string
	LayerNameSourceSetting         string
	AdditionalLayerInformation     *AdditionalLayerInformation

	isBlendClippedElements  bool
	isBlendInteriorElements bool
	isKnockout              bool
	ProtectionFlags         uint32
	Color                   *Color
	MetadataSetting         *MetadataSetting
	ReferencePoint          []float64
}

type ChannelInformation struct {
	Id          uint16
	Length      uint64
	Compression int16 //Compression. 0 = Raw Data, 1 = RLE compressed, 2 = ZIP without prediction, 3 = ZIP with prediction.
	Data        []byte
}

type LayerMaskORAdjustmentLayerData struct {
	DataSize            uint32
	EnclosingLayerMask  []*Rectangle
	DefaultColor        byte
	Flags               byte
	MaskParameters      byte // Only present if bit 4 of Flags set above.
	MaskParametersValue []byte
	Padding             uint16
	RealFlags           byte
	MaskBackground      byte
}

type LayerBlendingRangesData struct {
	LayerBlendingLength uint32

	SourceBlacks []int16
	SourceWhites []int16
	DestBlacks   []int16
	DestWhites   []int16
}

type AdditionalLayerInformation struct {
	Signature  string
	Key        string
	DataLength uint64
}

type EffectsLayerInfo struct {
	Version      uint16
	EffectsCount uint16
	Signature    string
}

type MetadataSetting struct {
	MetadataItemsCount     uint32
	DataSignature          []string
	DataKey                []string
	CopyOnSheetDuplication []byte
	Padding                [][]byte
	DataLength             []uint32
	Data                   [][]byte
}

/*****************************************************************************
 *
 * Image Data
 *
 *****************************************************************************
 */
type ImageData struct {
	CompressionMethod uint16
	Image             *image.RGBA
}
