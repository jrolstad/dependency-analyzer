package orchestration

import (
	"bufio"
	"github.com/jrolstad/dependency-analyzer/internal/models"
	"github.com/jrolstad/dependency-analyzer/internal/services"
	"os"
	"strings"
)

func GetDependencies(path string, filePattern string, fileService services.FileSearchService) ([]*models.DependencyNode, error) {
	files, err := fileService.Search(path, filePattern)
	if err != nil {
		return make([]*models.DependencyNode, 0), err
	}

	parsedDependencies := make([]*models.DependencyNode, 0)
	processingErrors := make([]error, 0)
	for _, file := range files {
		tree, err := parseFile(file)
		if err != nil {
			processingErrors = append(processingErrors, err)
			continue
		}

		parsedDependencies = append(parsedDependencies, tree...)
	}

	return parsedDependencies, nil
}

func parseFile(filePath string) ([]*models.DependencyNode, error) {
	fileContents, err := readFileContents(filePath)
	if err != nil {
		return make([]*models.DependencyNode, 0), nil
	}

	parsedData := &models.DependencyNode{
		Parent:    nil,
		FullName:  "",
		Namespace: "",
		Name:      "",
		Version:   "",
		Scope:     "",
		Children:  make(map[string]*models.DependencyNode),
	}

	if len(fileContents) == 1 {
		addDependencyDataToNode(fileContents[0], parsedData)
	} else {
		for _, line := range fileContents {
			// Do not process empty
			if strings.TrimSpace(line) == "" {
				continue
			}
			// First Line of the file
			if strings.HasPrefix(line, "digraph") {
				parentDependencyDataRaw := parseValueBetweenQuotes(line)
				addDependencyDataToNode(parentDependencyDataRaw, parsedData)
				parsedData.Children[parsedData.FullName] = parsedData
			}
			// Last Line of the file
			if strings.HasPrefix(line, "{") {
				continue
			}
			// Do the rest
			relationshipData := strings.Split(line, "->")
			if len(relationshipData) >= 2 {
				//parent := cleanDependencyName(parseValueBetweenQuotes(relationshipData[0]))
				//child := cleanDependencyName(parseValueBetweenQuotes(relationshipData[1]))

			}

		}
	}
	return []*models.DependencyNode{parsedData}, nil

}

func readFileContents(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func parseValueBetweenQuotes(s string) string {
	start := strings.Index(s, "\"")
	end := strings.Index(s[start+1:], "\"")
	if start >= 0 && end >= 0 {
		return s[start+1 : start+1+end]
	}
	return ""
}

func addDependencyDataToNode(rawData string, parsedData *models.DependencyNode) {
	splitData := strings.Split(rawData, ":")

	parsedData.FullName = cleanDependencyName(rawData)
	parsedData.Namespace = splitData[0]
	parsedData.Name = splitData[1]
	parsedData.Type = splitData[2]
	parsedData.Version = splitData[3]

	if len(splitData) > 4 {
		parsedData.Scope = splitData[4]
	}
}

func createNode(rawData string) *models.DependencyNode {
	splitData := strings.Split(rawData, ":")

	parsedData := &models.DependencyNode{}
	parsedData.FullName = cleanDependencyName(rawData)
	parsedData.Namespace = splitData[0]
	parsedData.Name = splitData[1]
	parsedData.Type = splitData[2]
	parsedData.Version = splitData[3]
	parsedData.Children = make(map[string]*models.DependencyNode)

	if len(splitData) > 4 {
		parsedData.Scope = splitData[4]
	}

	return parsedData
}

func cleanDependencyName(name string) string {
	splitData := strings.Split(name, " ")

	cleanValue := strings.Replace(splitData[0], "(", "", 1)
	cleanValue = strings.Replace(cleanValue, ")", "", 1)
	cleanValue = strings.TrimSpace(cleanValue)

	return cleanValue
}
