package main

import (
	"testing"
)

func TestInsertDependencyIntoList(t *testing.T) {
	rdl := createReverseDependencyList()

	countDependencies := len(rdl.Dependencies)
	if countDependencies != 2 {
		t.Error("Expected two dependencies; Got %v", countDependencies)
	}

	if rdl.Dependencies[0].PackageName != "glide" || rdl.Dependencies[1].PackageName != "golang" {
		t.Error("Expected glide and golang dependencies to be in list")
	}
}

func TestRemoveDependencyFromList(t *testing.T) {
	rdl := createReverseDependencyList()

	rdl.RemoveDependency(NewPackage("golang"))
	countDependencies := len(rdl.Dependencies)
	pkgName := rdl.Dependencies[0].PackageName

	if countDependencies != 1 || pkgName != "glide" {
		t.Errorf("Expected one dependency named glide; Got %v and %v",
			countDependencies, pkgName)
	}
}

func createReverseDependencyList() *ReverseDependencyList {
	rdl := NewReverseDependencyList()

	golang := NewPackage("golang")
	glide := NewPackage("glide")
	rdl.InsertNewDependency(golang)
	rdl.InsertNewDependency(glide)
	return rdl
}
