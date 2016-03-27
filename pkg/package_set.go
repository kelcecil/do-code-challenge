package pkg

import (
	"errors"
)

var (
	DEPENDENCY_NOT_AVAILABLE = errors.New("Dependency is not available.")
	REQUIRED_BY_OTHERS       = errors.New("Package is a dependency of other packages.")

	packages *PackageSet = NewPackageSet()
)

type PackageSet struct {
	Packages            map[string]*Package
	ReverseDependencies map[string]*ReverseDependencyList
}

func NewPackageSet() *PackageSet {
	return &PackageSet{
		Packages:            map[string]*Package{},
		ReverseDependencies: map[string]*ReverseDependencyList{},
	}
}

func (rs *PackageSet) FetchPackage(pkgName string) *Package {
	pack := rs.findPackage(pkgName)
	return pack
}

func (rs *PackageSet) findPackage(packageName string) *Package {
	pkg, ok := rs.Packages[packageName]
	if !ok {
		return nil
	}
	return pkg
}

func (rs *PackageSet) RemovePackage(pkgName string) error {
	checkPackage := rs.findPackage(pkgName)
	if checkPackage == nil {
		return nil
	}

	rdl, ok := rs.ReverseDependencies[pkgName]

	if !ok || rdl.IsDependedOn() {
		return REQUIRED_BY_OTHERS
	}

	pkg := rs.Packages[pkgName]

	for i := range pkg.Dependencies {
		dependency := pkg.Dependencies[i]
		rs.ReverseDependencies[dependency.PackageName].RemoveDependency(pkg)
	}

	delete(rs.Packages, pkgName)
	return nil
}

func (rs *PackageSet) InsertPackage(pkgName string, dependencies ...string) error {
	depPackages, ok := rs.FindRequiredDependencies(dependencies...)
	if !ok {
		return DEPENDENCY_NOT_AVAILABLE
	}

	if rs.findPackage(pkgName) == nil {
		rs.ReverseDependencies[pkgName] = NewReverseDependencyList()

		newPackage := NewPackage(pkgName, depPackages...)
		rs.Packages[pkgName] = newPackage

		for i := range dependencies {
			rs.ReverseDependencies[dependencies[i]].InsertNewDependency(newPackage)
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