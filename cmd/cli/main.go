package main

import (
	"flag"
	"fmt"
	"github.com/jrolstad/dependency-analyzer/internal/orchestration"
)

func main() {
	var path = ""
	var filePattern = ""

	flag.StringVar(&path, "path", ".", "Path to read files from")
	flag.StringVar(&filePattern, "filePattern", "dependencytree.dot", "Dependency file naming pattern")

	flag.Parse()

	dependencies, err := orchestration.GetDependencies(path, filePattern)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Dependencies:%v", len(dependencies))
}
