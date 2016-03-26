package main

import (
	"fmt"
	"sort"
)

type Message struct {
	Command             string
	PackageName         string
	PackageDependencies []string
	Response            chan string
}

func NewEmptyMessage() *Message {
	return &Message{
		Response: make(chan string),
	}
}

func NewMessage(cmd string, pkgName string, pkgDeps []string) *Message {
	return &Message{
		Command:             cmd,
		PackageName:         pkgName,
		PackageDependencies: pkgDeps,
		Response:            make(chan string),
	}
}

func (message *Message) String() string {
	if message == nil {
		return ""
	}
	return fmt.Sprintf("Command: %v, PackageName: %v, PackageDependencies %v\n",
		message.Command, message.PackageName, message.PackageDependencies)
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
