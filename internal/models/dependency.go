package models

type DependencyNode struct {
	Parent   *DependencyNode
	FullName string
	Name     string
	Version  string
	Children []*DependencyNode
}
