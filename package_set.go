package main

import (
	"errors"
	"sync"
)

var (
	DEPENDENCY_NOT_AVAILABLE = errors.New("Dependency is not available")
)

type PackageSet struct {
	Packages            map[string]*Package
	ReverseDependencies map[string][]*Package
	ReadWriteLock       *sync.RWMutex
}

func NewPackageSet() *PackageSet {
	return &PackageSet{
		Packages:            map[string]*Package{},
		ReverseDependencies: map[string][]*Package{},
		ReadWriteLock:       &sync.RWMutex{},
	}
}

func (rs *PackageSet) FetchPackages(PackageNames ...string) (Packages []*Package) {
	if rs == nil {
		return nil
	}

	rs.ReadWriteLock.RLock()
	defer rs.ReadWriteLock.RUnlock()

	for i := range PackageNames {
		Package := rs.findPackage(PackageNames[i])
		if Package != nil {
			Packages = append(Packages, Package)
		}
	}
	return Packages
}

func (rs *PackageSet) findPackage(packageName string) *Package {
	pkg, ok := rs.Packages[packageName]
	if !ok {
		return nil
	}
	return pkg
}

func (rs *PackageSet) InsertPackage(pkgName string, dependencies ...string) error {
	if rs == nil {
		return nil
	}

	rs.ReadWriteLock.Lock()
	defer rs.ReadWriteLock.Unlock()

	depPackages, ok := rs.FindRequiredDependencies(dependencies...)
	if !ok {
		return DEPENDENCY_NOT_AVAILABLE
	}

	if rs.findPackage(pkgName) == nil {
		rs.ReverseDependencies[pkgName] = make([]*Package, 0)

		newPackage := NewPackage(pkgName, depPackages...)
		rs.Packages[pkgName] = newPackage

		for i := range dependencies {
			rs.ReverseDependencies[dependencies[i]] =
				append(rs.ReverseDependencies[dependencies[i]], newPackage)
		}
	}
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
