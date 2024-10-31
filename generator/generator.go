package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	copyFile("thirdparty/stb/stb_image.h", "stb_image.h")
}

func copyFile(src, dst string) {
	_ = os.MkdirAll(filepath.Dir(dst), 0750)
	srcFile, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()

	_, err = srcFile.WriteTo(dstFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Copied file %s to %s\n", src, dst)
}
