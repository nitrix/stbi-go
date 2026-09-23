# stbi-go

A binding library of stb_image for Go.

## Usage

```go
import (
	"github.com/nitrix/stbi-go"
)

func example() error {
	img, err := stbi.Load("example.jpg")
	if err != nil {
		return err
	}

	// Do what you want with `img` here.
	// It's an `*image.NRGBA` (non-premultiplied alpha) with the pixel data in `.Pix` as usual.

	return nil
}
```

Every loader has a `Raw` variant that skips copying the pixels into Go memory. The returned slice is backed by C memory, so it must be released with `stbi.Free` and not used afterwards. This is useful when the pixels are handed off right away, for example uploaded to the GPU.

```go
pix, width, height, err := stbi.LoadRaw("example.jpg")
if err != nil {
	return err
}
defer stbi.Free(pix)
```

| Source | Image (Go memory) | Raw (C memory, call `Free`) |
|---|---|---|
| Path | `Load` | `LoadRaw` |
| `*os.File` | `LoadFile` | `LoadFileRaw` |
| `[]byte` | `LoadMemory` | `LoadMemoryRaw` |
| `io.Reader` | `LoadReader` | `LoadReaderRaw` |
| Path, as floats (HDR) | `Loadf` → `*stbi.NRGBAF32` | `LoadfRaw` → `[]float32` |

Pixels are always decoded to 4 channels (RGBA).

## Credits

See [this repo](https://github.com/nothings/stb) for the original C library.

## License

This is free and unencumbered software released into the public domain. See the [UNLICENSE](UNLICENSE) file for more details.