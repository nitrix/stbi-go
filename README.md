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
    // It's an `*image.RGBA` with the pixel data in `.Pix` as usual.

    return nil
}
```

There's also `LoadFile` to load from an `*os.File`, `LoadMemory` to load from a `[]byte` and `Loadf` which loads HDR images.

## Credits

See [this repo](https://github.com/nothings/stb) for the original C library.

## License

This is free and unencumbered software released into the public domain. See the [UNLICENSE](UNLICENSE) file for more details.