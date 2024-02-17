package services

import (
	"bufio"
	"fmt"
	"github.com/jrolstad/dependency-analyzer/internal/models"
	"os"
	"strings"
)

type DependencyParser interface {
	ParseFile(filePath string) ([]*models.DependencyNode, map[string]*models.DependencyNode, error)
}

func NewDependencyParser() DependencyParser {
	return &DependencyParserImpl{}
}

type DependencyParserImpl struct {
}

func (d *DependencyParserImpl) ParseFile(filePath string) ([]*models.DependencyNode, map[string]*models.DependencyNode, error) {
	fileContents, err := readFileContents(filePath)
	if err != nil {
		return make([]*models.DependencyNode, 0), make(map[string]*models.DependencyNode), nil
	}

	parsedData := createEmptyNode()

	allDependencies := make(map[string]*models.DependencyNode)
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
				allDependencies[parsedData.FullName] = parsedData

			}
			// Last Line of the file
			if strings.HasPrefix(line, "{") {
				continue
			}
			// Do the rest
			relationshipData := strings.Split(line, "->")
			if len(relationshipData) >= 2 {
				parentRaw := cleanDependencyName(parseValueBetweenQuotes(relationshipData[0]))
				parent := sanitizeDependencyName(parentRaw)
				childRaw := cleanDependencyName(parseValueBetweenQuotes(relationshipData[1]))
				child := sanitizeDependencyName(childRaw)

				if strings.HasSuffix(parentRaw, "test") || strings.HasSuffix(childRaw, "test") {
					continue
				}

				if allDependencies[parent] == nil {
					allDependencies[parent] = createNode(parentRaw)
				}

				if allDependencies[child] == nil {
					allDependencies[child] = createNode(childRaw)
				}

				parentNode := allDependencies[parent]

				if parentNode.Children[child] == nil {
					allDependencies[child].Parents[parentNode.FullName] = parentNode
					parentNode.Children[child] = allDependencies[child]
				}
			}

		}
	}
	return []*models.DependencyNode{parsedData}, allDependencies, nil

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

func createEmptyNode() *models.DependencyNode {
	return &models.DependencyNode{
		Parents:   make(map[string]*models.DependencyNode),
		FullName:  "",
		Namespace: "",
		Name:      "",
		Version:   "",
		Scope:     "",
		Children:  make(map[string]*models.DependencyNode),
	}
}

func addDependencyDataToNode(rawData string, parsedData *models.DependencyNode) {
	splitData := strings.Split(rawData, ":")

	parsedData.FullName = sanitizeDependencyName(cleanDependencyName(rawData))
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
	parsedData.Parents = make(map[string]*models.DependencyNode)
	parsedData.FullName = sanitizeDependencyName(cleanDependencyName(rawData))
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

func sanitizeDependencyName(name string) string {
	parts := strings.Split(name, ":")

	sanitizedName := fmt.Sprintf("%s:%s:%s:%s", parts[0], parts[1], parts[2], parts[3])

	return sanitizedName
}
