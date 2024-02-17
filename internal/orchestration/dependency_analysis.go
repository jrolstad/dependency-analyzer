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
		all, err := dependencyParser.ParseFile(file)
		if err != nil {
			processingErrors = append(processingErrors, err)
			continue
		}

		allParsedDependencies = append(allParsedDependencies, all)
	}

	return parsedDependencies, allParsedDependencies, nil
}

func IdentifyInScopeDependencies(dependencies []map[string]*models.DependencyNode) map[string]*models.DependencyNode {
	inScope := make(map[string]*models.DependencyNode)

	for _, item := range dependencies {
		for _, value := range item {
			if value.Parents != nil {
				for _, parent := range value.Parents {
					if strings.HasPrefix(parent.FullName, "com.oracle") &&
						!strings.HasPrefix(value.FullName, "com.oracle") &&
						!strings.HasPrefix(value.FullName, "javax.") &&
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

func IdentifyDependenciesNotReferencedByOthers(toAnalyze map[string]*models.DependencyNode) map[string]*models.DependencyNode {
	toRemove := getDependenciesReferencedByOthers(toAnalyze)

	results := make(map[string]*models.DependencyNode)
	for _, item := range toAnalyze {
		if toRemove[item.FullName] == nil {
			results[item.FullName] = item
		}
	}

	return results
}

func getDependenciesReferencedByOthers(toAnalyze map[string]*models.DependencyNode) map[string]*models.DependencyNode {
	toRemove := make(map[string]*models.DependencyNode)
	for _, item := range toAnalyze {
		for _, parent := range item.Parents {
			if !strings.HasPrefix(parent.FullName, "com.oracle") {
				toRemove[item.FullName] = item
			}
		}
	}
	return toRemove
}
