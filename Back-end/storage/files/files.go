package files

import (
	"errors"
	"fmt"
	"path/filepath"

	"read-adviser-bot/lib/e"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

var (
	ErrBadFileExtension = errors.New("bad file extension")
)

//Saving file
func IsValidExtension (filePath string) (error) {
	validExtensions := map[string]struct{}{
		".jpg": {},
		".jpeg": {},
		".png": {},
	}
	ext := filepath.Ext(filePath)
	fmt.Println("ext:", ext)

	if _, is := validExtensions[ext]; !is {
		return e.Wrap("Extension = " + ext, ErrBadFileExtension)
	}
	return nil
}
