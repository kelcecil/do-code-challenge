package main

import "testing"

func TestMessagesAreEqual(t *testing.T) {
	left := &Message{
		Command:             "INDEX",
		PackageName:         "glide",
		PackageDependencies: []string{"golang", "git", "hg"},
	}
	right := &Message{
		Command:             "INDEX",
		PackageName:         "glide",
		PackageDependencies: []string{"git", "hg", "golang"},
	}

	if !left.Equals(right) {
		t.Error("Messages that should be equal are not equal.")
	}
}

func TestCompareOrderedStringArrays(t *testing.T) {
	left := []string{"Golang", "Rust", "Crystal"}
	right := []string{"Golang", "Rust", "Crystal"}

	if !stringArraysAreEqual(left, right) {
		t.Errorf("Test for ordered string array equality fails.")
	}
}

func TestCompareUnorderedStringArrays(t *testing.T) {
	left := []string{"Golang", "Crystal", "Rust"}
	right := []string{"Golang", "Rust", "Crystal"}

	if !stringArraysAreEqual(left, right) {
		t.Errorf("Test for unordered string array equality fails.")
	}
}

func TestCompareDifferentStringArrays(t *testing.T) {
	left := []string{"Golang", "Crystal", "Rust"}
	right := []string{"Golang", "Rust", "Java"}

	if stringArraysAreEqual(left, right) {
		t.Errorf("Test for different items in string array fails.")
	}

	left = []string{"Golang", "Crystal", "Rust", "Java"}
	right = []string{"Golang", "Rust", "Java"}

	if stringArraysAreEqual(left, right) {
		t.Errorf("Arrays with unequal number of arguments passes.")
	}
}
