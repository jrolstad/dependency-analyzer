package services

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type FileSearchService interface {
	Search(root string, pattern string) ([]string, error)
}

func NewFileSearchService() FileSearchService {
	return &FileSearchServiceImpl{}
}

type FileSearchServiceImpl struct {
}

func (f FileSearchServiceImpl) Search(root string, pattern string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, visitDir(&files, pattern))
	if err != nil {
		return files, err
	}

	return files, err
}

func visitDir(files *[]string, pattern string) fs.WalkDirFunc {
	return func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		// Check if the entry is a file that matches the pattern
		if !d.IsDir() {
			matched, err := filepath.Match(pattern, filepath.Base(path))
			if err != nil {
				fmt.Println("Match error:", err)
				return err
			}

			if matched {
				*files = append(*files, path)
			}
		}

		return nil
	}
}
