package main

import (
	"flag"
	"fmt"
	"github.com/jrolstad/dependency-analyzer/internal/core"
	"github.com/jrolstad/dependency-analyzer/internal/orchestration"
	"github.com/jrolstad/dependency-analyzer/internal/services"
)

func main() {
	var path = ""
	var filePattern = ""

	flag.StringVar(&path, "path", ".", "Path to read files from")
	flag.StringVar(&filePattern, "filePattern", "dependencytree.dot", "Dependency file naming pattern")

	flag.Parse()

	fileService := services.NewFileSearchService()

	dependencies, err := orchestration.GetDependencies(path, filePattern, fileService)
	if err != nil {
		panic(err)
	}

	fmt.Print(core.MapToJson(dependencies))
}
