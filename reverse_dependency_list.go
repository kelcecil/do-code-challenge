package main

import (
	"sort"
)

type ReverseDependencyList struct {
	Dependencies []*Package
}

func NewReverseDependencyList() *ReverseDependencyList {
	return &ReverseDependencyList{
		Dependencies: []*Package{},
	}
}

func (rdl *ReverseDependencyList) IsDependedOn() bool {
	return rdl.Len() != 0
}

func (rdl *ReverseDependencyList) IsDependedOnBy(pkgName string) bool {
	i := sort.Search(rdl.Len(), func(i int) bool {
		return rdl.Dependencies[i].PackageName >= pkgName
	})
	if !(rdl.Len() == i) && rdl.Dependencies[i].PackageName == pkgName {
		return true
	}
	return false
}

func (rdl *ReverseDependencyList) InsertNewDependency(pkg *Package) bool {
	i := findIndice(rdl.Dependencies, pkg)

	if rdl.Len() == i ||
		rdl.Dependencies[i].PackageName != pkg.PackageName {
		rdl.Dependencies = insertIntoArray(rdl.Dependencies, pkg, i)
		return true
	}
	return false
}

func (rdl *ReverseDependencyList) RemoveDependency(pkg *Package) bool {
	i := findIndice(rdl.Dependencies, pkg)

	if !(rdl.Len() <= i) &&
		rdl.Dependencies[i].PackageName == pkg.PackageName {
		rdl.Dependencies = removeFromArray(rdl.Dependencies, i)
		return true
	}
	return false
}

func insertIntoArray(packages []*Package, pkg *Package, i int) []*Package {
	return append(packages[:i],
		append([]*Package{pkg}, packages[i:]...)...)
}

func removeFromArray(packages []*Package, i int) []*Package {
	return append(packages[:i], packages[i+1:]...)
}

func findIndice(packages []*Package, pkg *Package) int {
	return sort.Search(len(packages), func(i int) bool {
		return packages[i].PackageName >= pkg.PackageName
	})
}

func (rdl *ReverseDependencyList) Len() int {
	return len(rdl.Dependencies)
}

func (rdl *ReverseDependencyList) Swap(i, j int) {
	rdl.Dependencies[i], rdl.Dependencies[j] = rdl.Dependencies[j], rdl.Dependencies[i]
}

func (rdl *ReverseDependencyList) Less(i, j int) bool {
	return rdl.Dependencies[i].PackageName < rdl.Dependencies[j].PackageName
}
