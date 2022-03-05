package render

// screen resolution: 160x144
// internal resolution: 256x256

// FrameBuffer is typedef'd to the internal resolution, where the value at an element is the
// 2-bit grayscale value.  To support color, this should be a 3D uint8 array, where the last rank is 8-bit rgb
type FrameBuffer [256][256]uint8

const (
	// SDL colors
	WHITE = uint32(0xFFFFFFFF)
	DARK_GRAY = uint32(0xFF545454)
	LIGHT_GRAY = uint32(0xFFA8A8A8)
	BLACK = uint32(0xFF000000)
)


type PPU struct {
	FB FrameBuffer

	LCDC uint8 // FF40
}

// Init initializes the framebuffer
func (p *PPU) Init() {
	p.FB = FrameBuffer{}
}

