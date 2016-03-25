package main

import (
	"errors"
	"sort"
	"sync"
)

type Package struct {
	PackageName  string
	Dependencies []*Package
	DependedOn   []*Package
}

var (
	DEPENDENCY_NOT_AVAILABLE = errors.New("Dependency is not available")
)

func NewPackage(name string, packages ...*Package) *Package {
	return &Package{
		PackageName:  name,
		Dependencies: packages,
	}
}

type PackageSet struct {
	Packages      []*Package
	ReadWriteLock *sync.RWMutex
}

func NewPackageSet() *PackageSet {
	return &PackageSet{
		Packages:      []*Package{},
		ReadWriteLock: &sync.RWMutex{},
	}
}

func (rs *PackageSet) Items() []*Package {
	if rs == nil {
		return nil
	}

	if !sort.IsSorted(rs) {
		sort.Sort(rs)
	}
	return rs.Packages
}

func (rs *PackageSet) FetchPackages(PackageNames ...string) (Packages []*Package) {
	if rs == nil {
		return nil
	}

	rs.ReadWriteLock.RLock()
	for i := range PackageNames {
		Package := rs.findPackage(PackageNames[i])
		if Package != nil {
			Packages = append(Packages, Package)
		}
	}
	rs.ReadWriteLock.RUnlock()
	return Packages
}

func (rs *PackageSet) findPackage(PackageName string) *Package {
	PackagesInSet := rs.Items()
	candidateIndice := sort.Search(len(PackagesInSet), func(i int) bool {
		return PackagesInSet[i].PackageName >= PackageName
	})
	if candidateIndice >= len(PackagesInSet) {
		return nil
	}
	candidatePackage := PackagesInSet[candidateIndice]
	if candidatePackage.PackageName == PackageName {
		return candidatePackage
	}
	return nil
}

func (rs *PackageSet) InsertPackage(pkgName string, dependencies ...string) error {
	if rs == nil {
		return nil
	}

	rs.ReadWriteLock.Lock()
	depPackages, ok := rs.FindRequiredDependencies(dependencies...)
	if !ok {
		return DEPENDENCY_NOT_AVAILABLE
	}
	if rs.findPackage(pkgName) == nil {
		newPackage := NewPackage(pkgName, depPackages...)
		rs.Packages = append(rs.Packages, newPackage)
	}
	rs.ReadWriteLock.Unlock()
	return nil
}

func (rs *PackageSet) FindRequiredDependencies(dependencies ...string) (foundDependencies []*Package, noMissingDeps bool) {
	noMissingDeps = true
	for i := range dependencies {
		Package := rs.findPackage(dependencies[i])
		if Package != nil {
			foundDependencies = append(foundDependencies, Package)
		} else {
			noMissingDeps = false
		}
	}
	return foundDependencies, noMissingDeps
}

func (rs *PackageSet) Len() int {
	return len(rs.Packages)
}

func (rs *PackageSet) Swap(i, j int) {
	rs.Packages[i], rs.Packages[j] = rs.Packages[j], rs.Packages[i]
}

func (rs *PackageSet) Less(i, j int) bool {
	return rs.Packages[i].PackageName < rs.Packages[j].PackageName
}
