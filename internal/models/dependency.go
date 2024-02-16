package models

type DependencyNode struct {
	Parent    *DependencyNode
	FullName  string
	Name      string
	Namespace string
	Type      string
	Version   string
	Scope     string
	Children  map[string]*DependencyNode
}
