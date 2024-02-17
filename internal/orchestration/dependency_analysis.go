package orchestration

import (
	"github.com/jrolstad/dependency-analyzer/internal/models"
	"github.com/jrolstad/dependency-analyzer/internal/services"
	"strings"
)

func GetDependencies(path string, filePattern string, fileService services.FileSearchService, dependencyParser services.DependencyParser) ([]*models.DependencyNode, []map[string]*models.DependencyNode, error) {
	files, err := fileService.Search(path, filePattern)
	if err != nil {
		return make([]*models.DependencyNode, 0), make([]map[string]*models.DependencyNode, 0), err
	}

	parsedDependencies := make([]*models.DependencyNode, 0)
	allParsedDependencies := make([]map[string]*models.DependencyNode, 0)
	processingErrors := make([]error, 0)
	for _, file := range files {
		tree, all, err := dependencyParser.ParseFile(file)
		if err != nil {
			processingErrors = append(processingErrors, err)
			continue
		}

		parsedDependencies = append(parsedDependencies, tree...)
		allParsedDependencies = append(allParsedDependencies, all)
	}

	return parsedDependencies, allParsedDependencies, nil
}

func IdentifyInScopeIdentities(dependencies []map[string]*models.DependencyNode) map[string]*models.DependencyNode {
	inScope := make(map[string]*models.DependencyNode)

	for _, item := range dependencies {
		for _, value := range item {
			if value.Parents != nil {
				for _, parent := range value.Parents {
					if strings.HasPrefix(parent.FullName, "com.oracle") &&
						!strings.HasPrefix(value.FullName, "com.oracle") &&
						!strings.EqualFold(value.Scope, "test") &&
						!strings.EqualFold(value.Scope, "provided") {
						inScope[value.FullName] = value
					}
				}

			}
		}
	}

	return inScope
}

func IdentifyInScopeDependenciesNotReferencedByOthers(toAnalyze map[string]*models.DependencyNode, allDependencies []map[string]*models.DependencyNode) map[string]*models.DependencyNode {
	result := make(map[string]*models.DependencyNode)
	for _, item := range toAnalyze {
		for _, subset := range allDependencies {
			for _, all := range subset {
				if !strings.HasPrefix(all.FullName, "com.oracle") &&
					all.Children[item.FullName] == nil {
					result[item.FullName] = item
				}
			}
		}
	}

	return result
}
