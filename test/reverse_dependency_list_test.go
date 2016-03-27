package test

import (
	"github.com/kelcecil/do-code-challenge/pkg"
	"testing"
)

func TestInsertDependencyIntoList(t *testing.T) {
	rdl := createReverseDependencyList()

	countDependencies := len(rdl.Dependencies)
	if countDependencies != 2 {
		t.Errorf("Expected two dependencies; Got %v", countDependencies)
	}

	if rdl.Dependencies[0].PackageName != "glide" || rdl.Dependencies[1].PackageName != "golang" {
		t.Error("Expected glide and golang dependencies to be in list")
	}
}

func TestRemoveDependencyFromList(t *testing.T) {
	rdl := createReverseDependencyList()

	rdl.RemoveDependency(pkg.NewPackage("golang"))
	countDependencies := len(rdl.Dependencies)
	pkgName := rdl.Dependencies[0].PackageName

	if countDependencies != 1 || pkgName != "glide" {
		t.Errorf("Expected one dependency named glide; Got %v and %v",
			countDependencies, pkgName)
	}
}

func createReverseDependencyList() *pkg.ReverseDependencyList {
	rdl := pkg.NewReverseDependencyList()

	golang := pkg.NewPackage("golang")
	glide := pkg.NewPackage("glide")
	rdl.InsertNewDependency(golang)
	rdl.InsertNewDependency(glide)
	return rdl
}

func TestNoDuplicate(t *testing.T) {
	rdl := createReverseDependencyList()
	if rdl.InsertNewDependency(pkg.NewPackage("golang")) {
		t.Error("Duplicate dependency entry should not be allowed.")
	}
}

func TestIsDependedOnBy(t *testing.T) {
	rdl := createReverseDependencyList()
	if !rdl.IsDependedOnBy("golang") {
		t.Error("Golang should be a Reverse Dependency List.")
	}
	if rdl.IsDependedOnBy("scala") {
		t.Error("Scala should not be found in Reverse Dependency List")
	}
}

func TestReverseDependencyListSortFuncs(t *testing.T) {
	rdl := createReverseDependencyList()
	if rdl.Len() != 2 {
		t.Errorf("Len function is not correct. Want: 2, Got: %v", rdl.Len())
	}
	less := rdl.Less(0, 1)
	if !less {
		t.Errorf("Less function is not correct. Want: false, Got: %v", less)
	}
	rdl.Swap(0, 1)
	less = rdl.Less(0, 1)
	if less {
		t.Error("Swap function did not swap indices.")
	}
}
