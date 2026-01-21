package cmd

// Module directory constants
const (
	DirComponents = "components"
	DirBases      = "bases"
	DirProjects   = "projects"
)

// Module type constants
const (
	TypeComponent = "component"
	TypeBase      = "base"
	TypeProject   = "project"
)

// ModuleDirs contains all module directory names
var ModuleDirs = []string{DirComponents, DirBases, DirProjects}

// ModuleTypeOrder defines the sorting order for module types
var ModuleTypeOrder = map[string]int{
	TypeComponent: 1,
	TypeBase:      2,
	TypeProject:   3,
}

// ModuleInfo holds information about a discovered module
type ModuleInfo struct {
	Name    string
	Type    string
	Path    string
	Version string
}
