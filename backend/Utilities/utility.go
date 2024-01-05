package Utilities

import (
	"io"
	"os"
	"strings"
)

func MoveFile(src, dst string) error {
	extsplit := strings.Split(src, ".")[1]
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	out, err := os.Create(dst + "." + extsplit)
	if err != nil {
		in.Close()
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	in.Close()
	if err != nil {
		return err
	}

	err = out.Sync()
	if err != nil {
		return err
	}

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return err
	}
	return nil
}
