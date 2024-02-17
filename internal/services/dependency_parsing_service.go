package services

import (
	"bufio"
	"fmt"
	"github.com/jrolstad/dependency-analyzer/internal/models"
	"os"
	"strings"
)

type DependencyParser interface {
	ParseFile(filePath string) (map[string]*models.DependencyNode, error)
}

func NewDependencyParser() DependencyParser {
	return &DependencyParserImpl{}
}

type DependencyParserImpl struct {
}

func (d *DependencyParserImpl) ParseFile(filePath string) (map[string]*models.DependencyNode, error) {
	fileContents, err := readFileContents(filePath)
	if err != nil {
		return make(map[string]*models.DependencyNode), nil
	}

	dependenciesByName := make(map[string]*models.DependencyNode)

	if len(fileContents) > 1 {
		for _, line := range fileContents {
			// Do not process empty
			if strings.TrimSpace(line) == "" {
				continue
			}
			// First Line of the file
			if strings.HasPrefix(line, "digraph") {
				continue
			}
			// Last Line of the file
			if strings.HasPrefix(line, "{") {
				continue
			}
			// Do the rest
			relationshipData := strings.Split(line, "->")
			if len(relationshipData) >= 2 {
				parentRaw := cleanDependencyName(parseValueBetweenQuotes(relationshipData[0]))
				childRaw := cleanDependencyName(parseValueBetweenQuotes(relationshipData[1]))

				parentNode := createNode(parentRaw)
				childNode := createNode(childRaw)

				if d.isTestDependency(parentNode, childNode) {
					continue
				}

				if dependenciesByName[parentNode.FullName] == nil {
					dependenciesByName[parentNode.FullName] = parentNode
				}

				if dependenciesByName[childNode.FullName] == nil {
					dependenciesByName[childNode.FullName] = childNode
				}

				resolvedParentNode := dependenciesByName[parentNode.FullName]
				resolvedChildNode := dependenciesByName[childNode.FullName]

				resolvedParentNode.Children[resolvedChildNode.FullName] = resolvedChildNode
				resolvedChildNode.Parents[resolvedParentNode.FullName] = resolvedParentNode
			}

		}
	}
	return dependenciesByName, nil

}

func (d *DependencyParserImpl) isTestDependency(parentNode *models.DependencyNode, childNode *models.DependencyNode) bool {
	return strings.EqualFold(parentNode.Scope, "test") || strings.EqualFold(childNode.Scope, "test")
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
