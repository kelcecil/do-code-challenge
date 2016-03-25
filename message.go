package main

import "sort"

type Message struct {
	Command             string
	PackageName         string
	PackageDependencies []string
}

func NewMessage(cmd string, pkgName string, pkgDeps []string) *Message {
	return &Message{
		Command:             cmd,
		PackageName:         pkgName,
		PackageDependencies: pkgDeps,
	}
}

func (left *Message) Equals(right *Message) bool {
	if left == nil || right == nil {
		return false
	}

	stringsAreEqual := func(l string, r string) bool {
		return l == r
	}

	if !stringsAreEqual(left.Command, right.Command) {
		return false
	}

	if !stringsAreEqual(left.PackageName, right.PackageName) {
		return false
	}

	if !stringArraysAreEqual(left.PackageDependencies, right.PackageDependencies) {
		return false
	}

	return true
}

func stringArraysAreEqual(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}

	sort.Strings(right)

	for i := range left {
		rightIndice := sort.SearchStrings(right, left[i])
		if left[i] != right[rightIndice] {
			return false
		}
	}
	return true
}
