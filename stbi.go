package stbi

import (
	"errors"
	"image"
	"io"
	"os"
	"unsafe"
)

//go:generate go run ./generator

// #cgo LDFLAGS: -lm
// #include "stb_image.h"
import "C"

func Load(path string) (*image.RGBA, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	var x, y C.int
	data := C.stbi_load(cpath, &x, &y, nil, 4)
	if data == nil {
		msg := C.GoString(C.stbi_failure_reason())
		return nil, errors.New(msg)
	}
	defer C.stbi_image_free(unsafe.Pointer(data))

	return &image.RGBA{
		Pix:    C.GoBytes(unsafe.Pointer(data), y*x*4),
		Stride: 4,
		Rect:   image.Rect(0, 0, int(x), int(y)),
	}, nil
}

func Loadf(path string) (dt []float32, w int, h int, comp int, mfree func(), err error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	// I vaguely remember needing this?
	// C.stbi_set_flip_vertically_on_load(1)

	var tw, th, tcomp C.int
	data := C.stbi_loadf(cpath, &tw, &th, &tcomp, 0)
	if data == nil {
		msg := C.GoString(C.stbi_failure_reason())
		return nil, int(tw), int(th), int(tcomp), nil, errors.New(msg)
	}

	s := unsafe.Slice((*float32)(unsafe.Pointer(data)), int(tw*th*tcomp))

	return s, int(tw), int(th), int(tcomp), func() { C.stbi_image_free(unsafe.Pointer(data)) }, nil
}

func LoadFile(f *os.File) (*image.RGBA, error) {
	mode := C.CString("rb")
	defer C.free(unsafe.Pointer(mode))
	fp, err := C.fdopen(C.int(f.Fd()), mode)
	if err != nil {
		return nil, err
	}

	var x, y C.int
	data := C.stbi_load_from_file(fp, &x, &y, nil, 4)
	if data == nil {
		msg := C.GoString(C.stbi_failure_reason())
		return nil, errors.New(msg)
	}
	defer C.stbi_image_free(unsafe.Pointer(data))

	return &image.RGBA{
		Pix:    C.GoBytes(unsafe.Pointer(data), y*x*4),
		Stride: 4,
		Rect:   image.Rect(0, 0, int(x), int(y)),
	}, nil
}

func LoadMemory(b []byte) (*image.RGBA, error) {
	var x, y C.int
	mem := (*C.uchar)(unsafe.Pointer(&b[0]))
	data := C.stbi_load_from_memory(mem, C.int(len(b)), &x, &y, nil, 4)
	if data == nil {
		msg := C.GoString(C.stbi_failure_reason())
		return nil, errors.New(msg)
	}
	defer C.stbi_image_free(unsafe.Pointer(data))

	return &image.RGBA{
		Pix:    C.GoBytes(unsafe.Pointer(data), y*x*4),
		Stride: 4,
		Rect:   image.Rect(0, 0, int(x), int(y)),
	}, nil
}

func LoadReader(r io.Reader) (*image.RGBA, error) {
	if f, ok := r.(*os.File); ok {
		return LoadFile(f)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return LoadMemory(b)
}
