package goPsdLib

type effectsLayer interface {
	readEffectsLayer(doc *document)
}

type CommonStateInfo struct {
	NextThreeItems uint32
	Version        uint32
	Visible        bool
	Unused         uint16
}

type DropAndInnerShadowInfo struct {
	// Size of the remaining items: 41 or 51 (depending on version)
	RemainingItems uint32
	// Version: 0 ( Photoshop 5.0) or 2 ( Photoshop 5.5)
	Version uint32
	// 	Blur value in pixels
	PixelsBlurValue uint32
	// 	Intensity as a percent
	Intensity uint32
	// Angle in degrees
	AngleDegrees uint32
	// Distance in pixels
	Distance uint32
	// Color: 2 bytes for space followed by 4 * 2 byte color component
	Color []byte
	// Blend mode: 4 bytes for signature and 4 bytes for key
	BlendMode []byte
	// Effect enabled
	Enabled bool
	// Use this angle in all of the layer effects
	isUseAllLayer bool
	// Opacity as a percent
	Opacity byte
	// Native color: 2 bytes for space followed by 4 * 2 byte color component
	NativeColor []byte
}

func (commonStateInfo *CommonStateInfo) readEffectsLayer(doc *document) {
	doc.readUint32(&commonStateInfo.NextThreeItems)
	doc.readUint32(&commonStateInfo.Version)
	doc.readBool(&commonStateInfo.Visible)
	doc.readUint16(&commonStateInfo.Unused)

}

func (dropAndInnerShadowInfo *DropAndInnerShadowInfo) readEffectsLayer(doc *document) {
	doc.readUint32(&dropAndInnerShadowInfo.RemainingItems)
	doc.readUint32(&dropAndInnerShadowInfo.Version)
	doc.readUint32(&dropAndInnerShadowInfo.PixelsBlurValue)
	doc.readUint32(&dropAndInnerShadowInfo.Intensity)
	doc.readUint32(&dropAndInnerShadowInfo.AngleDegrees)
	doc.readUint32(&dropAndInnerShadowInfo.Distance)
	dropAndInnerShadowInfo.Color = doc.Next(10)
	dropAndInnerShadowInfo.BlendMode = doc.Next(8)
	doc.readBool(&dropAndInnerShadowInfo.Enabled)
	doc.readBool(&dropAndInnerShadowInfo.isUseAllLayer)

	dropAndInnerShadowInfo.Opacity, _ = doc.ReadByte()

	dropAndInnerShadowInfo.NativeColor = doc.Next(10)
}
