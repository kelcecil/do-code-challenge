package test

import (
	"github.com/kelcecil/do-code-challenge/message"
	"testing"
)

func TestMessagesAreEqual(t *testing.T) {
	left := message.NewMessage("INDEX", "glide", []string{"golang", "git", "hg"})
	right := message.NewMessage("INDEX", "glide", []string{"git", "hg", "golang"})

	if !left.Equals(right) {
		t.Error("Messages that should be equal are not equal.")
	}
}

func TestMessageCommandsAreDifferent(t *testing.T) {
	left := message.NewMessage("WINDEX", "glide", []string{"golang", "git", "hg"})
	right := message.NewMessage("INDEX", "glide", []string{"git", "hg", "golang"})

	if left.Equals(right) {
		t.Error("Message commands are different and should not be equal.")
	}
}

func TestMessagePackagesAreDifferent(t *testing.T) {
	left := message.NewMessage("INDEX", "slide", []string{"golang", "git", "hg"})
	right := message.NewMessage("INDEX", "glide", []string{"git", "hg", "golang"})

	if left.Equals(right) {
		t.Error("Message packages are different and should not be equal.")
	}
}

func TestPassingNilMessageShouldReturnFalse(t *testing.T) {
	left := message.NewMessage("WINDEX", "slide", []string{"golang", "git", "hg"})
	if left.Equals(nil) {
		t.Error("Nil pointer for package should not be equal")
	}
}

func TestMessageDependenciesAreDifferent(t *testing.T) {
	left := message.NewMessage("INDEX", "glide", []string{"golang", "git", "hg"})
	right := message.NewMessage("INDEX", "glide", []string{"git", "svn", "golang"})

	if left.Equals(right) {
		t.Error("Dependencies are different and should not match.")
	}
}

func TestCompareOrderedStringArrays(t *testing.T) {
	left := []string{"Golang", "Rust", "Crystal"}
	right := []string{"Golang", "Rust", "Crystal"}

	if !message.StringArraysAreEqual(left, right) {
		t.Errorf("Test for ordered string array equality fails.")
	}
}

func TestCompareUnorderedStringArrays(t *testing.T) {
	left := []string{"Golang", "Crystal", "Rust"}
	right := []string{"Golang", "Rust", "Crystal"}

	if !message.StringArraysAreEqual(left, right) {
		t.Errorf("Test for unordered string array equality fails.")
	}
}

func TestCompareDifferentStringArrays(t *testing.T) {
	left := []string{"Golang", "Crystal", "Rust"}
	right := []string{"Golang", "Rust", "Java"}

	if message.StringArraysAreEqual(left, right) {
		t.Errorf("Test for different items in string array fails.")
	}

	left = []string{"Golang", "Crystal", "Rust", "Java"}
	right = []string{"Golang", "Rust", "Java"}

	if message.StringArraysAreEqual(left, right) {
		t.Errorf("Arrays with unequal number of arguments passes.")
	}
}
