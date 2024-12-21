package util

import "path/filepath"

func ToAbsPath(path string) string{
	p, err := filepath.Abs(path)
	if err != nil{
		return path
	}
	return p
}