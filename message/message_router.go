package message

import (
	"github.com/kelcecil/do-code-challenge/pkg"
)

var (
	packages *pkg.PackageSet = pkg.NewPackageSet()
)

func MessageRouter(messages <-chan *Message) {
	for {
		msg := <-messages
		if msg.Command == "INDEX" {
			index(msg)
		} else if msg.Command == "QUERY" {
			query(msg)
		} else if msg.Command == "REMOVE" {
			remove(msg)
		} else {
			unknown(msg)
		}
	}
}

type CommandFunc func(*Message)

func createCommandResponseFunc(method func(msg *Message) interface{}) CommandFunc {
	return func(msg *Message) {
		output := method(msg)
		if output == nil {
			msg.Response <- "OK"
		} else {
			msg.Response <- "FAIL"
		}
	}
}

var index = createCommandResponseFunc(func(msg *Message) interface{} {
	return packages.InsertPackage(msg.PackageName, msg.PackageDependencies...)
})

var remove = createCommandResponseFunc(func(msg *Message) interface{} {
	return packages.RemovePackage(msg.PackageName)
})

func query(msg *Message) {
	pkg := packages.FetchPackage(msg.PackageName)
	if pkg != nil {
		msg.Response <- "OK"
	} else {
		msg.Response <- "FAIL"
	}
}

func unknown(msg *Message) {
	msg.Response <- "ERROR"
}
