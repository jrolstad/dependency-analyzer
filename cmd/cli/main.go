package main

import (
	"flag"
	"fmt"
	"github.com/jrolstad/dependency-analyzer/internal/models"
	"github.com/jrolstad/dependency-analyzer/internal/orchestration"
	"github.com/jrolstad/dependency-analyzer/internal/services"
	"sort"
)

func main() {
	var path = ""
	var filePattern = ""

	flag.StringVar(&path, "path", ".", "Path to read files from")
	flag.StringVar(&filePattern, "filePattern", "dependencytree.dot", "Dependency file naming pattern")

	flag.Parse()

	fileService := services.NewFileSearchService()
	parser := services.NewDependencyParser()

	_, allDependencies, err := orchestration.GetDependencies(path, filePattern, fileService, parser)
	if err != nil {
		panic(err)
	}

	inScope := orchestration.IdentifyInScopeDependencies(allDependencies)
	inScopeNotReferenced := orchestration.IdentifyDependenciesNotReferencedByOthers(inScope)

	showData(inScopeNotReferenced)
}

func showData(data map[string]*models.DependencyNode) {

	keys := make([]string, 0, len(data))

	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Println(data[key].FullName)
	}
}
