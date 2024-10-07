package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	err := copyFile("thirdparty/stb/stb_image.h", "stb_image.h")
	if err != nil {
		panic(err)
	}
}

func copyFile(src, dst string) error {
	_ = os.MkdirAll(filepath.Dir(dst), 0750)
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = srcFile.WriteTo(dstFile)
	if err != nil {
		return err
	}

	fmt.Printf("Copied file %s to %s\n", src, dst)

	return nil
}
