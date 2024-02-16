package orchestration

import (
	"fmt"
	"github.com/jrolstad/dependency-analyzer/internal/models"
	"github.com/jrolstad/dependency-analyzer/internal/services"
)

func GetDependencies(path string, filePattern string, fileService services.FileSearchService) ([]*models.DependencyNode, error) {
	files, err := fileService.Search(path, filePattern)
	if err != nil {
		return make([]*models.DependencyNode, 0), err
	}

	for _, file := range files {
		fmt.Println(file)
	}

	return make([]*models.DependencyNode, 0), nil
}
