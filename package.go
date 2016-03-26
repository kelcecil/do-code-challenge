package main

type Package struct {
	PackageName  string
	Dependencies []*Package
	DependedOn   []*Package
}

func NewPackage(name string, packages ...*Package) *Package {
	return &Package{
		PackageName:  name,
		Dependencies: packages,
	}
}
