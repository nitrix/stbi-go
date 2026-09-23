package stbi

import (
	"errors"
	"image"
	"image/color"
	"io"
	"math"
	"os"
	"slices"
	"unsafe"
)

//go:generate go run ./generator

// #cgo CFLAGS: -O3
// #cgo LDFLAGS: -lm
// #include "stb_image.h"
import "C"

// Pixels are always decoded to 4 channels (RGBA), non-premultiplied alpha.
const channels = 4

// Load decodes the image file at path into Go memory.
func Load(path string) (*image.NRGBA, error) {
	return toNRGBA(LoadRaw(path))
}

// LoadRaw decodes the image file at path without copying the pixels into Go memory.
// The returned slice is backed by C memory and must be released with Free.
func LoadRaw(path string) (pix []byte, width, height int, err error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	var x, y C.int
	data := C.stbi_load(cpath, &x, &y, nil, channels)

	return wrap[byte](unsafe.Pointer(data), x, y)
}

// Loadf decodes the image file at path as floats into Go memory, keeping HDR values above 1.
func Loadf(path string) (*NRGBAF32, error) {
	return toNRGBAF32(LoadfRaw(path))
}

// LoadfRaw decodes the image file at path as floats without copying them into Go memory.
// The returned slice is backed by C memory and must be released with Free.
func LoadfRaw(path string) (pix []float32, width, height int, err error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	var x, y C.int
	data := C.stbi_loadf(cpath, &x, &y, nil, channels)

	return wrap[float32](unsafe.Pointer(data), x, y)
}

// LoadFile decodes the image read from f into Go memory.
func LoadFile(f *os.File) (*image.NRGBA, error) {
	return toNRGBA(LoadFileRaw(f))
}

// LoadFileRaw decodes the image read from f without copying the pixels into Go memory.
// The returned slice is backed by C memory and must be released with Free.
func LoadFileRaw(f *os.File) (pix []byte, width, height int, err error) {
	return LoadReaderRaw(f)
}

// LoadMemory decodes the encoded image in b into Go memory.
func LoadMemory(b []byte) (*image.NRGBA, error) {
	return toNRGBA(LoadMemoryRaw(b))
}

// LoadMemoryRaw decodes the encoded image in b without copying the pixels into Go memory.
// The returned slice is backed by C memory and must be released with Free.
func LoadMemoryRaw(b []byte) (pix []byte, width, height int, err error) {
	if len(b) == 0 {
		return nil, 0, 0, errors.New("empty image data")
	}

	if len(b) > math.MaxInt32 {
		return nil, 0, 0, errors.New("image data too large")
	}

	var x, y C.int
	mem := (*C.stbi_uc)(unsafe.Pointer(&b[0]))
	data := C.stbi_load_from_memory(mem, C.int(len(b)), &x, &y, nil, channels)

	return wrap[byte](unsafe.Pointer(data), x, y)
}

// LoadReader decodes the image read from r into Go memory.
func LoadReader(r io.Reader) (*image.NRGBA, error) {
	return toNRGBA(LoadReaderRaw(r))
}

// LoadReaderRaw decodes the image read from r without copying the pixels into Go memory.
// The returned slice is backed by C memory and must be released with Free.
func LoadReaderRaw(r io.Reader) (pix []byte, width, height int, err error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, 0, 0, err
	}

	return LoadMemoryRaw(b)
}

// Free releases pixels returned by one of the Raw functions. The slice must be
// exactly the one returned, and must not be used afterwards. Freeing nil is a no-op.
func Free[T byte | float32](pix []T) {
	if len(pix) == 0 {
		return
	}

	C.stbi_image_free(unsafe.Pointer(&pix[0]))
}

// NRGBAF32 is an in-memory image of non-premultiplied float32 RGBA pixels.
// Values are not clamped, so HDR images can exceed 1.
type NRGBAF32 struct {
	Pix    []float32
	Stride int // In float32s, not bytes.
	Rect   image.Rectangle
}

func (p *NRGBAF32) ColorModel() color.Model {
	return color.NRGBA64Model
}

func (p *NRGBAF32) Bounds() image.Rectangle {
	return p.Rect
}

// At converts to color.NRGBA64, clamping values to [0, 1].
func (p *NRGBAF32) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.NRGBA64{}
	}

	i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*channels
	s := p.Pix[i : i+channels : i+channels]

	return color.NRGBA64{
		R: toUint16(s[0]),
		G: toUint16(s[1]),
		B: toUint16(s[2]),
		A: toUint16(s[3]),
	}
}

func toUint16(v float32) uint16 {
	return uint16(min(max(v, 0), 1)*0xffff + 0.5)
}

// wrap turns pixel data returned by stb_image into a slice over the same C memory.
func wrap[T byte | float32](data unsafe.Pointer, x, y C.int) ([]T, int, int, error) {
	if data == nil {
		msg := C.GoString(C.stbi_failure_reason())
		return nil, 0, 0, errors.New(msg)
	}

	width, height := int(x), int(y)

	return unsafe.Slice((*T)(data), width*height*channels), width, height, nil
}

// toNRGBA copies raw pixels into Go memory and frees them.
func toNRGBA(pix []byte, width, height int, err error) (*image.NRGBA, error) {
	if err != nil {
		return nil, err
	}
	defer Free(pix)

	return &image.NRGBA{
		Pix:    slices.Clone(pix),
		Stride: width * channels,
		Rect:   image.Rect(0, 0, width, height),
	}, nil
}

// toNRGBAF32 copies raw pixels into Go memory and frees them.
func toNRGBAF32(pix []float32, width, height int, err error) (*NRGBAF32, error) {
	if err != nil {
		return nil, err
	}
	defer Free(pix)

	return &NRGBAF32{
		Pix:    slices.Clone(pix),
		Stride: width * channels,
		Rect:   image.Rect(0, 0, width, height),
	}, nil
}
