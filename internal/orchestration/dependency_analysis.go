package orchestration

import (
	"github.com/jrolstad/dependency-analyzer/internal/models"
	"github.com/jrolstad/dependency-analyzer/internal/services"
)

func GetDependencies(path string, filePattern string, fileService services.FileSearchService, dependencyParser services.DependencyParser) ([]*models.DependencyNode, error) {
	files, err := fileService.Search(path, filePattern)
	if err != nil {
		return make([]*models.DependencyNode, 0), err
	}

	parsedDependencies := make([]*models.DependencyNode, 0)
	processingErrors := make([]error, 0)
	for _, file := range files {
		tree, _, err := dependencyParser.ParseFile(file)
		if err != nil {
			processingErrors = append(processingErrors, err)
			continue
		}

		parsedDependencies = append(parsedDependencies, tree...)
	}

	return parsedDependencies, nil
}
