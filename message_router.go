package main

import "log"

func MessageRouter(messages <-chan *Message) {
	for {
		select {
		case msg := <-messages:
			log.Print(msg)
			if msg.Command == "INDEX" {
				index(msg)
			} else if msg.Command == "QUERY" {
				query(msg)
			} else if msg.Command == "REMOVE" {
				remove(msg)
			} else {
				unknown(msg)
			}
		default:
		}
	}
}

func index(msg *Message) {
	err := packages.InsertPackage(msg.PackageName, msg.PackageDependencies...)
	if err == nil {
		msg.Response <- "OK"
	} else {
		msg.Response <- "FAIL"
	}
}

func query(msg *Message) {
	pkg := packages.FetchPackage(msg.PackageName)
	if pkg != nil {
		msg.Response <- "OK"
	} else {
		msg.Response <- "FAIL"
	}
}

func remove(msg *Message) {
	err := packages.RemovePackage(msg.PackageName)
	if err == nil {
		msg.Response <- "OK"
	} else {
		msg.Response <- "FAIL"
	}
}

func unknown(msg *Message) {
	msg.Response <- "ERROR"
}
