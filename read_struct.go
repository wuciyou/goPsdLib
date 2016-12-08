package goPsdLib

import (
	"log"
)

func (rectangle *Rectangle) readStructure(doc *document) {
	doc.readUint32(&rectangle.Top)
	doc.readUint32(&rectangle.Left)
	doc.readUint32(&rectangle.Bottom)
	doc.readUint32(&rectangle.Right)

	rectangle.Y = rectangle.Top
	rectangle.X = rectangle.Left
	rectangle.Width = rectangle.Right - rectangle.Left
	rectangle.Height = rectangle.Bottom - rectangle.Top

}

func (descriptorStructure *DescriptorStructure) readStructure(doc *document) {
	doc.readUnicodeString(&descriptorStructure.Name)
	doc.readDynamicString(&descriptorStructure.ClassId)
	doc.readUint32(&descriptorStructure.ItemsNumber)
	itemsNumber := int(descriptorStructure.ItemsNumber)
	descriptorStructure.Items = make(map[string]*Descriptor)
	for item_index := 0; item_index < itemsNumber; item_index++ {
		descriptor := &Descriptor{}
		descriptor.readStructure(doc)
		descriptorStructure.Items[descriptor.Type] = descriptor
		// log.Printf("descriptor:%+v \n ", descriptor)
	}
}

func (descriptor *Descriptor) readStructure(doc *document) {

	if !descriptor.isSkipKey {
		doc.readDynamicString(&descriptor.Key)
	}

	descriptor.Type = string(doc.Next(4))
	// log.Printf("doc  descriptor：%+v, data:%+v ", descriptor, doc.Bytes())
	switch descriptor.Type {
	case "obj ": // Reference
		reference := &ReferenceStructureEntity{}
		reference.readStructure(doc)
		descriptor.Value = reference
	case "Objc", "GlbO": // Descriptor      // GlobalObject same as Descriptor
		descriptorStructure := &DescriptorStructure{}
		descriptorStructure.readStructure(doc)
		descriptor.Value = descriptorStructure

	case "VlLs": // List
		list := &ListStructureEntity{}
		list.readStructure(doc)
		descriptor.Value = list
	case "doub": // Double
	case "UntF": // Unit float
	case "TEXT": // String
		stringStructure := &StringStructureEntity{}
		stringStructure.readStructure(doc)
		descriptor.Value = stringStructure
	case "enum": // Enumerated
		enumeratedDescriptor := &EnumeratedDescriptorEntity{}
		enumeratedDescriptor.readStructure(doc)
		descriptor.Value = enumeratedDescriptor
	case "long": // Integer
		long := &LongEntity{}
		long.readStructure(doc)
		descriptor.Value = long

	case "comp": // Large Integer
	case "bool": // Boolean
		boolE := &BooleanStructureEntity{}
		boolE.readStructure(doc)
		descriptor.Value = boolE
	case "type": // Class
	case "GlbC": // Class
	case "alis": // Alias
	case "tdta": // Raw Data
	default:
		log.Panicf("Unknown OSType key [%s] in entity [%s]", descriptor.Key, descriptor.Type)
	}

	if descriptor.Value == nil {
		log.Panicf("Unknown implement key [%s] in entity [%s]", descriptor.Key, descriptor.Type)
	}

	// log.Printf("descriptor:%+v \n ", descriptor)
}

func (list *ListStructureEntity) readStructure(doc *document) {
	doc.readUint32(&list.ItemsNumber)
	itemNumber := int(list.ItemsNumber)
	list.DescriptorList = make([]*Descriptor, itemNumber)
	for i := 0; i < itemNumber; i++ {
		descriptor := &Descriptor{isSkipKey: true}
		descriptor.readStructure(doc)
		list.DescriptorList[i] = descriptor
	}

}

func (reference *ReferenceStructureEntity) readStructure(doc *document) {
	var itemsNumber uint32
	doc.readUint32(&itemsNumber)
	itemsNum := int(itemsNumber)

	reference.entitys = make(map[string]StructureEntity)
	for i := 0; i < itemsNum; i++ {
		Type := string(doc.Next(4))
		switch Type {
		case "prop": // Property
			property := &PropertyStructureEntity{}
			property.readStructure(doc)
			reference.entitys[Type] = property
		case "Clss": // Class
			class := &ClassStructureEntity{}
			class.readStructure(doc)
			reference.entitys[Type] = class
		case "Enmr": // Enumerated Reference
			enumerated := &EnumeratedReferenceEntity{}
			enumerated.readStructure(doc)
			reference.entitys[Type] = enumerated
		case "rele": // Offset
			offset := &OffsetStructureEntity{}
			offset.readStructure(doc)
			reference.entitys[Type] = offset
		case "Idnt": // Identifier
		case "indx": // Index
		case "name": //Name
		}

	}
}

func (property *PropertyStructureEntity) readStructure(doc *document) {
	doc.readUnicodeString(&property.Name)
	doc.readDynamicString(&property.ClassID)
	doc.readDynamicString(&property.KeyId)
}

func (class *ClassStructureEntity) readStructure(doc *document) {
	doc.readUnicodeString(&class.Name)
	doc.readDynamicString(&class.ClassID)
}
func (enumeratedReference *EnumeratedReferenceEntity) readStructure(doc *document) {
	doc.readUnicodeString(&enumeratedReference.Name)
	doc.readDynamicString(&enumeratedReference.ClassID)
	doc.readDynamicString(&enumeratedReference.TypeID)
	doc.readDynamicString(&enumeratedReference.Enum)
}

func (offset *OffsetStructureEntity) readStructure(doc *document) {
	doc.readUnicodeString(&offset.Name)
	doc.readDynamicString(&offset.ClassID)
	doc.readUint32(&offset.Value)
}
func (long *LongEntity) readStructure(doc *document) {
	var longData uint32
	doc.readUint32(&longData)
	long.Value = uint64(longData)
}

func (enumerated *EnumeratedDescriptorEntity) readStructure(doc *document) {
	doc.readDynamicString(&enumerated.Type)
	doc.readDynamicString(&enumerated.Enum)
}

func (str *StringStructureEntity) readStructure(doc *document) {
	doc.readUnicodeString(&str.Value)
}

func (boolE *BooleanStructureEntity) readStructure(doc *document) {
	doc.readBool(&boolE.isTrue)
}

/****************************************************************
*
* Layer and Mask Information Section
*
*****************************************************************/
func (layer *LayerInfo) readStructure(doc *document) {

	if doc.isPSB {
		doc.readUint64(&layer.Length)
	} else {
		var layer_length uint32
		doc.readUint32(&layer_length)
		layer.Length = uint64(layer_length)
	}

	doc.readUint16(&layer.LayerCount)

	layerCount := int(layer.LayerCount)

	for i := 0; i < layerCount; i++ {
		layerRecords := &LayerRecords{}
		layerRecords.readStructure(doc)
		layer.LayerRecords = append(layer.LayerRecords, layerRecords)
	}

	for _, layerRecords := range layer.LayerRecords {
		width := int(layerRecords.Rectangle.Width)
		height := int(layerRecords.Rectangle.Height)

		for _, channelInformation := range layerRecords.ChannelInformations {

			doc.readInt16(&channelInformation.Compression)

			switch channelInformation.Compression {
			case 0: // Raw Data
				channelInformation.Data = doc.Next(width * height)
			case 1: // RLE compressed
				var result []byte
				scanLines := make([]int16, height)
				for i := range scanLines {
					doc.readInt16(&scanLines[i])
				}

				for i := range scanLines {
					dataTemp := doc.Next(int(scanLines[i]))
					line := UnpackRLEBits(dataTemp, width)
					result = append(result, line...)
				}
				channelInformation.Data = result
			// case 2: // ZIP without prediction
			// case 3: // ZIP with prediction
			default:
				log.Panicf("[Layer: %s] Unknown [compression:%d ] method of channel [id: %d]", layerRecords.LayerName, channelInformation.Compression, channelInformation.Id)
			}
		}
	}

}

func (layerRecords *LayerRecords) readStructure(doc *document) {

	layerRecords.Rectangle = doc.readRectangle()

	doc.readUint16(&layerRecords.ChannelsNumber)

	for c := 0; c < int(layerRecords.ChannelsNumber); c++ {
		channelInformation := &ChannelInformation{}
		doc.readUint16(&channelInformation.Id)
		if doc.isPSB {
			doc.readUint64(&channelInformation.Length)
		} else {
			var cLengU32 uint32
			doc.readUint32(&cLengU32)
			channelInformation.Length = uint64(cLengU32)
		}

		layerRecords.ChannelInformations = append(layerRecords.ChannelInformations, channelInformation)

	}

	layerRecords.BlendModeSignature = string(doc.Next(4))

	if layerRecords.BlendModeSignature != "8BIM" {
		log.Panicf("Wrong blend mode signature of layer [#%s].", layerRecords.BlendModeSignature)
	}

	layerRecords.BlendModeKey = string(doc.Next(4))

	layerRecords.Opacity, _ = doc.ReadByte()
	layerRecords.Clipping, _ = doc.ReadByte()
	layerRecords.Flags, _ = doc.ReadByte()
	layerRecords.Filler, _ = doc.ReadByte()
	doc.readUint32(&layerRecords.ExtraDataFieldLength)

	extraDataEntBuf := NewDocumentFromByte(doc.Next(int(layerRecords.ExtraDataFieldLength)))
	extraDataEntBuf.isPSB = doc.isPSB

	// Layer mask data [http://www.adobe.com/devnet-apps/photoshop/fileformatashtml/PhotoshopFileFormats.htm#50577409_22582]
	layerRecords.LayerMaskORAdjustmentLayerData = &LayerMaskORAdjustmentLayerData{}
	layerRecords.LayerMaskORAdjustmentLayerData.readStructure(extraDataEntBuf)

	// Layer blending ranges data [http://www.adobe.com/devnet-apps/photoshop/fileformatashtml/PhotoshopFileFormats.htm#50577409_21332]
	layerRecords.LayerBlendingRangesData = &LayerBlendingRangesData{}
	layerRecords.LayerBlendingRangesData.readStructure(extraDataEntBuf)

	extraDataEntBuf.readPascalString(&layerRecords.LayerName)

	nameLength := len(layerRecords.LayerName) + 1
	if nameLength%4 != 0 {
		skip := 4 - nameLength%4
		extraDataEntBuf.Next(skip)
	}

	for extraDataEntBuf.Len() > 0 {

		additionalLayerInformation := &AdditionalLayerInformation{}

		additionalLayerInformation.Signature = string(extraDataEntBuf.Next(4))
		additionalLayerInformation.Key = string(extraDataEntBuf.Next(4))

		if extraDataEntBuf.isPSB && StringValueIs(additionalLayerInformation.Key, "LMsk", "Lr16", "Lr32", "Layr", "Mt16", "Mt32", "Mtrn", "Alph", "FMsk", "lnk2", "FEid", "FXid", "PxSD") {
			extraDataEntBuf.readUint64(&additionalLayerInformation.DataLength)
		} else {
			var dataLength uint32
			extraDataEntBuf.readUint32(&dataLength)
			additionalLayerInformation.DataLength = uint64(dataLength)
		}
		additionalLength := int(additionalLayerInformation.DataLength + 1 & ^0x01)
		switch additionalLayerInformation.Key {
		case "luni":
			extraDataEntBuf.readUnicodeString(&layerRecords.LayerName)
		case "lnsr":
			layerRecords.LayerNameSourceSetting = string(extraDataEntBuf.Next(4))
		case "lyid":
			extraDataEntBuf.readUint32(&layerRecords.ID)
		case "clbl":
			extraDataEntBuf.readBool(&layerRecords.isBlendClippedElements)
			// Padding
			extraDataEntBuf.Next(3)
		case "infx":
			extraDataEntBuf.readBool(&layerRecords.isBlendInteriorElements)
			// Padding
			extraDataEntBuf.Next(3)
		case "knko":
			extraDataEntBuf.readBool(&layerRecords.isKnockout)
			extraDataEntBuf.Next(3)
		case "lspf":
			extraDataEntBuf.readUint32(&layerRecords.ProtectionFlags)
		case "lclr":
			layerRecords.Color = &Color{}
			layerRecords.Color.readStructure(extraDataEntBuf)
		case "shmd":
			layerRecords.MetadataSetting = &MetadataSetting{}
			layerRecords.MetadataSetting.readStructure(extraDataEntBuf)
		case "fxrp":
			referencePoint := make([]float64, 2)
			extraDataEntBuf.readInt64(&referencePoint[0])
			extraDataEntBuf.readInt64(&referencePoint[1])
			layerRecords.ReferencePoint = referencePoint
		default:
			log.Printf("key[%s]  skip:%d", additionalLayerInformation.Key, additionalLength)
			extraDataEntBuf.Next(additionalLength)
			// log.Panicf("unknown key[%s] \n", additionalLayerInformation.Key)
		}

	}
}

func (layerMask *LayerMaskORAdjustmentLayerData) readStructure(doc *document) {

	doc.readUint32(&layerMask.DataSize)

	layerMaskSize := int(layerMask.DataSize)

	layerMaskDataBuf := NewDocumentFromByte(doc.Next(layerMaskSize))

	if int(layerMask.DataSize) > 0 {

		layerMask.EnclosingLayerMask = append(layerMask.EnclosingLayerMask, layerMaskDataBuf.readRectangle())
		layerMask.DefaultColor, _ = layerMaskDataBuf.ReadByte()
		layerMask.Flags, _ = layerMaskDataBuf.ReadByte()

		if layerMask.Flags == 4 {
			layerMask.MaskParameters, _ = layerMaskDataBuf.ReadByte()
			switch layerMask.MaskParameters {
			case 0, 2:
				layerMask.MaskParametersValue = layerMaskDataBuf.Next(1)
			case 1, 3:
				layerMask.MaskParametersValue = layerMaskDataBuf.Next(8)
			}
		}

		if int(layerMask.DataSize) == 20 {
			layerMaskDataBuf.readUint16(&layerMask.Padding)
		} else {

			layerMask.RealFlags, _ = layerMaskDataBuf.ReadByte()
			layerMask.MaskBackground, _ = layerMaskDataBuf.ReadByte()

			layerMask.EnclosingLayerMask = append(layerMask.EnclosingLayerMask, layerMaskDataBuf.readRectangle())
		}

	}

}

func (layerBlending *LayerBlendingRangesData) readStructure(doc *document) {
	doc.readUint32(&layerBlending.LayerBlendingLength)

	layerBlendingCount := int(layerBlending.LayerBlendingLength / 8)

	layerBlending.SourceBlacks = make([]int16, layerBlendingCount)
	layerBlending.SourceWhites = make([]int16, layerBlendingCount)
	layerBlending.DestBlacks = make([]int16, layerBlendingCount)
	layerBlending.DestWhites = make([]int16, layerBlendingCount)

	for i := 0; i < layerBlendingCount; i++ {
		doc.readInt16(&layerBlending.SourceBlacks[i])
		doc.readInt16(&layerBlending.SourceWhites[i])
		doc.readInt16(&layerBlending.DestBlacks[i])
		doc.readInt16(&layerBlending.DestWhites[i])
	}
}

func (additionalLayerInformation *AdditionalLayerInformation) readStructure(doc *document) {

	log.Printf("additionalLayerInformation:%+v", additionalLayerInformation)
}

func (effectsLayerInfo *EffectsLayerInfo) readStructure(doc *document) {
	// Version: 0
	doc.readUint16(&effectsLayerInfo.Version)
	doc.readUint16(&effectsLayerInfo.EffectsCount)

	effectsCount := int(effectsLayerInfo.EffectsCount)
	for i := 0; i < effectsCount; i++ {
		effectsLayerInfo.Signature = string(doc.Next(4))
		if effectsLayerInfo.Signature != "8BIM" {
			log.Panicf("Wrong blend mode signature of layer [#%s].", effectsLayerInfo.Signature)
		}
		key := string(doc.Next(4))
		log.Printf("Key :%s \n ", key)
	}

}

func (c *Color) readStructure(doc *document) {
	doc.readUint16(&c.Red)
	doc.readUint16(&c.Green)
	doc.readUint16(&c.Blue)
	doc.readUint16(&c.Alpha)
}

func (metadatta *MetadataSetting) readStructure(doc *document) {

	doc.readUint32(&metadatta.MetadataItemsCount)
	var metadataItemsCount = int(metadatta.MetadataItemsCount)

	metadatta.DataSignature = make([]string, metadataItemsCount)
	metadatta.DataKey = make([]string, metadataItemsCount)
	metadatta.CopyOnSheetDuplication = make([]byte, metadataItemsCount)
	metadatta.Padding = make([][]byte, metadataItemsCount)
	metadatta.DataLength = make([]uint32, metadataItemsCount)
	metadatta.Data = make([][]byte, metadataItemsCount)
	for i := 0; i < metadataItemsCount; i++ {
		metadatta.DataSignature[i] = string(doc.Next(4))
		metadatta.DataKey[i] = string(doc.Next(4))
		metadatta.CopyOnSheetDuplication[i], _ = doc.ReadByte()
		metadatta.Padding[i] = doc.Next(3)

		doc.readUint32(&metadatta.DataLength[i])
		metadatta.Data[i] = doc.Next(int(metadatta.DataLength[i]))
	}
}
