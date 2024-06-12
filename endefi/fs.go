package endefi

import (
	"os"
	"path/filepath"
)

func ListDir(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, f := range files {
		if !f.IsDir() {
			names = append(names, path+"/"+f.Name())
		}
	}
	return names, nil
}

func WalkDir(path string) ([]string, error) {
	var names []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			names = append(names, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return names, nil
}
