package main

import "testing"

func TestSimplePackageFetch(t *testing.T) {
	packageSet := createNewPackageSetWithData()

	fetchedpackage := packageSet.FetchPackages("homebrew")
	if len(fetchedpackage) != 1 || fetchedpackage[0].PackageName != "homebrew" {
		t.Error("Fetching a single package failed.")
	}

	fetchedpackage = packageSet.FetchPackages("golang", "homebrew")
	if len(fetchedpackage) != 2 || fetchedpackage[0].PackageName != "golang" {
		t.Errorf("Fetching multiple Packages at once failed.")
	}
}

func TestInsertDuplicateNewPackage(t *testing.T) {
	packageSet := createNewPackageSetWithData()
	result := packageSet.InsertPackage("golang")
	if result != nil {
		t.Error("Duplicate Packages should not be possible.")
	}
}

func TestInsertNewPackageWithKnownGoodDependencies(t *testing.T) {
	packageSet := createNewPackageSetWithData()
	err := packageSet.InsertPackage("glide", "golang", "homebrew")
	if err != nil {
		t.Errorf("Error encountered when inserting with known good deps; Message: %v", err.Error())
	}
	newPackages := packageSet.FetchPackages("glide")

	packageCount := len(newPackages)
	if packageCount != 1 {
		t.Errorf("Expected one package returned; Got %v", packageCount)
	}

	packageName := newPackages[0].PackageName
	if packageName != "glide" {
		t.Errorf("Expected package to be named glide; Got %v", packageName)
	}

	lengthOfDeps := len(newPackages[0].Dependencies)
	if lengthOfDeps != 2 {
		t.Errorf("Expected two dependencies to be returned; Got %v", lengthOfDeps)
	}

	dependencyOneName := newPackages[0].Dependencies[0].PackageName
	dependencyTwoName := newPackages[0].Dependencies[1].PackageName
	depNamesMatch := dependencyOneName == "golang" &&
		dependencyTwoName == "homebrew"

	if !depNamesMatch {
		t.Errorf("Expected golang and homebrew for dependency names; Got %v and %v",
			dependencyOneName, dependencyTwoName)
	}
}

func TestInsertNewPackageWithKnownBadDependencies(t *testing.T) {
	packageSet := createNewPackageSetWithData()
	err := packageSet.InsertPackage("glide", "golang", "left-pad")
	if err == nil {
		t.Errorf("Inserting with known bad deps should not have happened")
	}
	if err != DEPENDENCY_NOT_AVAILABLE {
		t.Errorf("Expected: DEPENDENCY_NOT_AVAILABLE error; Got: %v", err)
	}
}

func TestFindKnownInsertedDependencies(t *testing.T) {
	packageSet := createNewPackageSetWithData()

	Packages, ok := packageSet.FindRequiredDependencies("golo", "homebrew")
	if len(Packages) != 2 || !ok {
		t.Errorf("Expected to find 2 dependencies and ok = true; Received %v dependencies and ok = %v", len(Packages), ok)
	}
	expectedDepOne := Packages[0].PackageName
	expectedDepTwo := Packages[1].PackageName
	if expectedDepOne != "golo" || expectedDepTwo != "homebrew" {
		t.Errorf("Expected golo and homebrew deps; Received %v and %v", expectedDepOne, expectedDepTwo)
	}
}

func TestFindKnownNotInsertedDependencies(t *testing.T) {
	packageSet := createNewPackageSetWithData()

	Packages, ok := packageSet.FindRequiredDependencies("java", "fish")
	packageCount := len(Packages)
	if packageCount != 0 || ok {
		t.Errorf("Expected to find 0 dependencies and ok = false; Received %v dependencies and ok = %v", packageCount, ok)
	}
}

func TestFindMixOfDependencies(t *testing.T) {
	packageSet := createNewPackageSetWithData()

	Packages, ok := packageSet.FindRequiredDependencies("golo", "fish")
	packageCount := len(Packages)
	if packageCount != 1 || ok {
		t.Errorf("Expected to find 1 dependencies and ok = false; Received %v dependencies and ok = %v", packageCount, ok)
	}
	expectedDep := Packages[0].PackageName
	if expectedDep != "golo" {
		t.Errorf("Expected golo deps; Received %v", expectedDep)
	}
}

func TestFetchNonExistentPackage(t *testing.T) {
	packageSet := createNewPackageSetWithData()

	fetchedPackage := packageSet.FetchPackages("harveyRabbit")
	if len(fetchedPackage) != 0 {
		t.Errorf("A package was fetched that should not exist.")
	}
}

func TestReverseDependencyTracking(t *testing.T) {
	packageSet := createNewPackageSetWithData()
	err := packageSet.InsertPackage("glide", "golang")
	if err != nil {
		t.Error("Inserting package failed.")
	}
	dependents := packageSet.ReverseDependencies["golang"]
	countOfDependents := len(dependents)
	if countOfDependents != 1 && dependents[0].PackageName != "glide" {
		t.Errorf("Expected reverse dependencies to have 1 and to be glide; Got %v and %v",
			countOfDependents, dependents[0].PackageName)
	}
}

func BenchmarkInsertPackages(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createNewPackageSetWithData()
	}
}

func createNewPackageSetWithData() *PackageSet {
	Packages := []string{"homebrew", "golang", "golo", "sdl"}

	packageSet := NewPackageSet()
	for i := range Packages {
		packageSet.InsertPackage(Packages[i])
	}
	return packageSet
}
