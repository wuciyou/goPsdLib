package goPsdLib

import (
// "log"
)

func imageResourceIDs1050(resource_data_buf *document, identifier uint16) *SlicesResourceFormat {
	slicesResourceFormat := &SlicesResourceFormat{}
	slicesHeader := &SlicesHeader{}

	resource_data_buf.readUint32(&slicesHeader.Version)

	switch slicesHeader.Version {
	case 7, 8:
		// resource_data_buf.readUint32(&slicesHeader.DescriptorVersion)
	case 6:
		resource_data_buf.readUint32(&slicesHeader.Top)
		resource_data_buf.readUint32(&slicesHeader.Leftt)
		resource_data_buf.readUint32(&slicesHeader.Bottom)
		resource_data_buf.readUint32(&slicesHeader.Right)
		resource_data_buf.readUnicodeString(&slicesHeader.GroupName)
		resource_data_buf.readUint32(&slicesHeader.SlicesNumber)

		slicesResourceFormat.Header = slicesHeader

		var blocks []*SlicesResourceBlock
		blocks = make([]*SlicesResourceBlock, slicesHeader.SlicesNumber)
		slicesNumber := int(slicesHeader.SlicesNumber)
		for i := 0; i < slicesNumber; i++ {
			block := &SlicesResourceBlock{}

			resource_data_buf.readUint32(&block.ID)
			resource_data_buf.readUint32(&block.GroupId)
			resource_data_buf.readUint32(&block.Origin)
			if block.Origin == 1 {
				resource_data_buf.readUint32(&block.AssociatedLayerID)
			}
			resource_data_buf.readUnicodeString(&block.Name)
			resource_data_buf.readUint32(&block.Type)
			resource_data_buf.readUint32(&block.Left)
			resource_data_buf.readUint32(&block.Top)
			resource_data_buf.readUint32(&block.Right)
			resource_data_buf.readUint32(&block.Bottom)
			resource_data_buf.readUnicodeString(&block.URL)
			resource_data_buf.readUnicodeString(&block.Target)
			resource_data_buf.readUnicodeString(&block.Message)
			resource_data_buf.readUnicodeString(&block.AltTag)
			resource_data_buf.readBool(&block.CellTextIsHTML)
			resource_data_buf.readUnicodeString(&block.CellText)

			resource_data_buf.readUint32(&block.HorizontalAlignment)
			resource_data_buf.readUint32(&block.VerticalAlignment)
			block.AlphaColor, _ = resource_data_buf.ReadByte()
			block.Red, _ = resource_data_buf.ReadByte()
			block.Green, _ = resource_data_buf.ReadByte()
			block.Blue, _ = resource_data_buf.ReadByte()

			resource_data_buf.readUint32(&block.DescriptorVersion)

			descriptorStructure := &DescriptorStructure{}

			descriptorStructure.readStructure(resource_data_buf)

			block.descriptorStructure = descriptorStructure

			blocks[i] = block

		}

		slicesResourceFormat.Blocks = blocks

	}

	return slicesResourceFormat
}
