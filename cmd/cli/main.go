package main

import (
	"flag"
	"fmt"
	"github.com/jrolstad/dependency-analyzer/internal/models"
	"github.com/jrolstad/dependency-analyzer/internal/orchestration"
	"github.com/jrolstad/dependency-analyzer/internal/services"
	"sort"
	"strings"
)

func main() {
	var path = ""
	var filePattern = ""
	var includedParentsRaw = ""
	var excludedDependenciesRaw = ""
	var excludedScopesRaw = ""
	var mode = ""

	flag.StringVar(&path, "path", ".", "Path to read files from")
	flag.StringVar(&filePattern, "filePattern", "dependencytree.dot", "Dependency file naming pattern")
	flag.StringVar(&includedParentsRaw, "includedParents", "", "Parent dependency naming pattern")
	flag.StringVar(&excludedDependenciesRaw, "excludedDependencies", "", "Dependency naming pattern to exclude")
	flag.StringVar(&excludedScopesRaw, "excludedScopes", "test,provided", "Dependency scopes to exclude")
	flag.StringVar(&mode, "mode", "notreferenced", "Analysis mode.  Valid values are all and notreferenced")
	flag.Parse()

	fileService := services.NewFileSearchService()
	parser := services.NewDependencyParser()

	_, allDependencies, err := orchestration.GetDependencies(path, filePattern, fileService, parser)
	if err != nil {
		panic(err)
	}

	includedParents := strings.Split(includedParentsRaw, ",")
	excludedDependencies := strings.Split(excludedDependenciesRaw, ",")
	excludedScopes := strings.Split(excludedScopesRaw, ",")

	inScope := orchestration.IdentifyInScopeDependencies(allDependencies, includedParents, excludedDependencies, excludedScopes)
	if strings.EqualFold(mode, "notreferenced") {
		inScopeNotReferenced := orchestration.IdentifyDependenciesNotReferencedByOthers(inScope, includedParents)
		showData(inScopeNotReferenced)
	} else {
		showData(inScope)
	}

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
