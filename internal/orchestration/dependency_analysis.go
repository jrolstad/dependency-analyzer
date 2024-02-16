package orchestration

import "github.com/jrolstad/dependency-analyzer/internal/models"

func GetDependencies(path string, filePattern string) ([]*models.DependencyNode, error) {
	return make([]*models.DependencyNode, 0), nil
}
