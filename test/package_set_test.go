package test

import (
	"github.com/kelcecil/do-code-challenge/pkg"
	"testing"
)

func TestSimplePackageFetch(t *testing.T) {
	packageSet := createNewPackageSetWithData()

	fetchedPackage := packageSet.FetchPackage("homebrew")
	if fetchedPackage.PackageName != "homebrew" {
		t.Error("Fetching a single package failed.")
	}
}

func TestInsertDuplicateNewPackage(t *testing.T) {
	packageSet := createNewPackageSetWithData()
	result := packageSet.InsertPackage("golang")
	if result != nil {
		t.Error("Duplicate Packages should not be possible.")
	}
}

func TestPackageRemovalWithNoDependents(t *testing.T) {
	packageSet := createNewPackageSetWithData()
	pkg := packageSet.FetchPackage("sdl")
	if pkg == nil {
		t.Error("Expected one package when asking for specific pkg before delete")
	}
	err := packageSet.RemovePackage("sdl")
	if err != nil {
		t.Errorf("Failure to remove package with no dependents; Message: %v", err)
	}
	pkg = packageSet.FetchPackage("sdl")
	if pkg != nil {
		t.Error("Expected zero packages when asking for specific pkg")
	}
}

func TestPackageRemovalWithNonexistentPackage(t *testing.T) {
	packageSet := createNewPackageSetWithData()
	err := packageSet.RemovePackage("harveyRabbit")
	if err != nil {
		t.Error("Removal of nonexistent pkg should have returned no error.")
	}
}

func TestPackageRemovalWithDependentsFails(t *testing.T) {
	packageSet := createNewPackageSetWithData()
	err := packageSet.RemovePackage("golang")
	if err != pkg.REQUIRED_BY_OTHERS {
		t.Errorf("Expected REQUIRED_BY_OTHERS error; Got %v", err)
	}
}

func TestPackageRemovalWithPreviousDependents(t *testing.T) {
	packageSet := createNewPackageSetWithData()
	err := packageSet.RemovePackage("glide")
	if err != nil {
		t.Errorf("Failed to remove dependent package; Error: %v", err)
	}
	err = packageSet.RemovePackage("golang")
	if err != nil {
		t.Errorf("Failed to remove package with former dependent; Error: %v", err)
	}
}

func TestInsertNewPackageWithKnownGoodDependencies(t *testing.T) {
	packageSet := createNewPackageSetWithData()
	pkg := packageSet.FetchPackage("glide")

	if pkg == nil {
		t.Errorf("Expected one package returned; Got zero")
	}

	if pkg.PackageName != "glide" {
		t.Errorf("Expected package to be named glide.")
	}

	lengthOfDeps := len(pkg.Dependencies)
	if lengthOfDeps != 2 {
		t.Errorf("Expected two dependencies to be returned; Got %v", lengthOfDeps)
	}

	dependencyOneName := pkg.Dependencies[0].PackageName
	dependencyTwoName := pkg.Dependencies[1].PackageName
	depNamesMatch := dependencyOneName == "golang" &&
		dependencyTwoName == "homebrew"

	if !depNamesMatch {
		t.Errorf("Expected golang and homebrew for dependency names; Got %v and %v",
			dependencyOneName, dependencyTwoName)
	}
}

func TestInsertNewPackageWithKnownBadDependencies(t *testing.T) {
	packageSet := createNewPackageSetWithData()
	err := packageSet.InsertPackage("do-example", "golang", "left-pad")
	if err == nil {
		t.Errorf("Inserting with known bad deps should not have happened")
	}
	if err != pkg.DEPENDENCY_NOT_AVAILABLE {
		t.Errorf("Expected: DEPENDENCY_NOT_AVAILABLE error; Got: %v", err)
	}
}

func TestFindKnownInsertedDependencies(t *testing.T) {
	packageSet := createNewPackageSetWithData()

	packages, ok := packageSet.FindRequiredDependencies("golo", "homebrew")
	if len(packages) != 2 || !ok {
		t.Errorf("Expected to find 2 dependencies and ok = true; Received %v dependencies and ok = %v", len(packages), ok)
	}
	expectedDepOne := packages[0].PackageName
	expectedDepTwo := packages[1].PackageName
	if expectedDepOne != "golo" || expectedDepTwo != "homebrew" {
		t.Errorf("Expected golo and homebrew deps; Received %v and %v", expectedDepOne, expectedDepTwo)
	}
}

func TestFindKnownNotInsertedDependencies(t *testing.T) {
	packageSet := createNewPackageSetWithData()

	packages, ok := packageSet.FindRequiredDependencies("java", "fish")
	packageCount := len(packages)
	if packageCount != 0 || ok {
		t.Errorf("Expected to find 0 dependencies and ok = false; Received %v dependencies and ok = %v", packageCount, ok)
	}
}

func TestFindMixOfDependencies(t *testing.T) {
	packageSet := createNewPackageSetWithData()

	packages, ok := packageSet.FindRequiredDependencies("golo", "fish")
	packageCount := len(packages)
	if packageCount != 1 || ok {
		t.Errorf("Expected to find 1 dependencies and ok = false; Received %v dependencies and ok = %v", packageCount, ok)
	}
	expectedDep := packages[0].PackageName
	if expectedDep != "golo" {
		t.Errorf("Expected golo deps; Received %v", expectedDep)
	}
}

func TestFetchNonExistentPackage(t *testing.T) {
	packageSet := createNewPackageSetWithData()

	fetchedPackage := packageSet.FetchPackage("harveyRabbit")
	if fetchedPackage != nil {
		t.Errorf("A package was fetched that should not exist.")
	}
}

func TestReverseDependencyTracking(t *testing.T) {
	packageSet := createNewPackageSetWithData()
	rdl := packageSet.ReverseDependencies["golang"]

	if !rdl.IsDependedOn() {
		t.Errorf("Expected package to be depended on")
	}
	if !rdl.IsDependedOnBy("glide") {
		t.Error("Expected package to be depended on by example.")
	}
}

func BenchmarkInsertPackages(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createNewPackageSetWithData()
	}
}

func BenchmarkFetchPackage(b *testing.B) {
	packageSet := createNewPackageSetWithData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		packageSet.FetchPackage("homebrew")
		packageSet.FetchPackage("golang")
		packageSet.FetchPackage("golo")
		packageSet.FetchPackage("sdl")
	}
}

func createNewPackageSetWithData() *pkg.PackageSet {
	packages := []string{"homebrew", "golang", "golo", "sdl"}

	packageSet := pkg.NewPackageSet()
	for i := range packages {
		packageSet.InsertPackage(packages[i])
	}

	packageSet.InsertPackage("glide", "golang", "homebrew")
	return packageSet
}
