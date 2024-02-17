package orchestration

import (
	"github.com/jrolstad/dependency-analyzer/internal/models"
	"github.com/jrolstad/dependency-analyzer/internal/services"
	"slices"
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

func IdentifyInScopeDependencies(dependencies []map[string]*models.DependencyNode,
	includedParents []string,
	excludedDependencies []string,
	excludedScopes []string) map[string]*models.DependencyNode {
	inScope := make(map[string]*models.DependencyNode)

	for _, item := range dependencies {
		for _, value := range item {
			if value.Parents != nil {
				for _, parent := range value.Parents {
					if hasItemThatStartsWith(includedParents, parent.FullName) &&
						!hasItemThatStartsWith(excludedDependencies, value.FullName) &&
						!slices.Contains(excludedScopes, value.Scope) {
						inScope[value.FullName] = value
					}
				}

			}
		}
	}

	return inScope
}

func hasItemThatStartsWith(allowed []string, value string) bool {
	for _, item := range allowed {
		if strings.HasPrefix(value, item) {
			return true
		}
	}

	return false
}

func IdentifyDependenciesNotReferencedByOthers(toAnalyze map[string]*models.DependencyNode, includedParents []string) map[string]*models.DependencyNode {
	toRemove := getDependenciesReferencedByOthers(toAnalyze, includedParents)

	results := make(map[string]*models.DependencyNode)
	for _, item := range toAnalyze {
		if toRemove[item.FullName] == nil {
			results[item.FullName] = item
		}
	}

	return results
}

func getDependenciesReferencedByOthers(toAnalyze map[string]*models.DependencyNode,
	includedParents []string) map[string]*models.DependencyNode {
	toRemove := make(map[string]*models.DependencyNode)
	for _, item := range toAnalyze {
		for _, parent := range item.Parents {
			if !hasItemThatStartsWith(includedParents, parent.FullName) {
				toRemove[item.FullName] = item
			}
		}
	}
	return toRemove
}
