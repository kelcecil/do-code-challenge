# do-code-challenge Design Document

## Overview

### Objective

Create a TCP server that acts as a "package" server. The server should known three commands.
* INDEX - Adds a package to the server.
* QUERY - Checks to see if a package is on the server.
* REMOVE - Removes package from the server.

Communication occurs using a newline (\n) terminated protocol as the following:

`<command>|<packageName>|<dependency1>,<dependency2>,...\n`

Both command and package name are required. Dependencies are optional and used only in the case of the INDEX command. Notice the line is terminated with newline.

There are three possible responses to commands: OK, FAIL, and ERROR. OK and FAIL are returned for different reasons for each command:
* INDEX
  * OK is returned if the package can be indexed or is already indexed.
  * FAIL is returned if the package can't be indexed because some of it's dependencies have yet to be installed.
* REMOVE
  * OK is returned if the package is removed from the package index or wasn't in the package index.
  * FAIL is returned if the package cannot be removed because another package depends on it.
* QUERY
  * OK is returned if the package is indexed.
  * FAIL is returned if the package isn't indexed.

ERROR is returned if the command isn't recognized or there is a problem with the request.

All responses are terminated with the newline character (e.g. `OK\n`, `FAIL\n`, `ERROR\n`).

### Requirements

Go 1.6 or newer, make

### Building and Running

Build and run included tests by running `make`. There are also separate make targets for building (`make build`) and running the included tests (`make test`).

Note that this program *must* be built *outside* of the GOPATH to allow for a multipackage program (like this one) to use relative import paths. Relative paths are not allowed inside the GOPATH to avoid ambiguity, but a remote import path would have personally identifiable information that would violate one of the submissions requirements. More information can be found at: https://golang.org/cmd/go/#hdr-Relative_import_paths

A Dockerfile is also provided to build and run tests on the indexer and create a container from the result binary:

`docker build -t indexer .`

If you'd like a randomly selected high port based on the Dockerfile's EXPOSE instruction:
`docker run -d -P indexer`

If you'd like to select the exposed port yourself:
`docker run -d -p <yourPreferredPort>:8080 indexer`

## Design

### Packages
The indexer is divided into several packages:
* message
  * This is code that deals with the convenience of creation and the handling of messages (think the newline delimited protocol.)
* parser
  * This package handles parsing the text protocol and creating messages from the reading.
* pkg
  * This package contains structures and operations that deal with packages.
  * `package_set.go` specifically holds methods related to working with the package store.
  * `reverse_dependency_list.go` holds a struct and functions that aid in determining if dependencies are associated with packages.
* server
  * This packages handles the TCP server.
* test
  * This packages holds all unit and integration tests for the indexer.

### Accepting Connections

The application begins by creating a TCP server using Go's `net` library and beginning to listen for connections. We set a periodic timeout on `listener.Accept()` to check for signals such as SIGTERM (which Docker will send on a first attempt to stop a container) so that we can handle a server stop gracefully by *attempting* to allow existing connections to finish; Docker will send a SIGKILL after a configurable amount of time following the SIGTERM signal if the process hasn't terminated and thus will interrupt remaining connections. New connections are spawned into a Go routine to be handled by `HandleConnection`.

### Parsing Message

The first step after establishing the connection is to read input from the client and parse. This is handled by the MessageReader in `parser.go`. The MessageReader reads 4096 bytes from the client at a time (as opposed to using `ReadString('\n')` to read until a newline terminator) to allow for us to inspect and determine if a client is just sending nonsense to attempt to overflow or cause a panic; This overflow check has not been added as a maximum message size was not defined in the problem description but could easily be added if one was provided.

The input from the client is added to a buffer until a newline (`\n`) is found. The buffer is then sent to the `ParseMessage` function in `parser.go` to be parsed. byte.buffer's `ReadString()` function is used extensively here since we now have a buffer of known size (and to keep the parsing code simple.) Situations where tokens are missing unexpectedly return an `error` stating such. A successful parsing returns a pointer to a `Message` that is returned to be sent down a buffered channel to our next topic.

### Message Router

Our parsed message is sent via buffered channel to the `MessageRouter`. The MessageRouter's responsibility is to direct the messages for commands to the correct function for interacting with the `PackageSet`. The MessageRouter likewise takes the result from these functions and translates to what the web server will return to the client (`OK`,`FAIL`, and `ERROR`).

### Interacting with the Package Set

The Package Set is the structure in which the state of the index is kept. The structure consists of two standard Go maps. One map uses the package name as a string as a key and Package objects as values. Package structs contains the package name and an array of pointers to Package objects representing the dependencies of the package. The second map keeps a structure called a Reverse Dependency List (`reverse_dependency_list.go`) which is an array of packages with helpful functions to assess packages that depend on a certain package.

Maps were chosen for storage over arrays using binary search (using Go's `sort.Search()` and implementing the `Sort` interface) because of the reduced average case complexity for inserts (amortized O(1) vs O(log n), search (amortized O(1) vs O(log n), and equal storage complexity (O(n) for both.) Arrays are still used for `ReverseDependencyList` and others as array performance using binary search can be better for a small number of items (see: http://www.darkcoding.net/software/go-slice-search-vs-map-lookup/).

The `MessageRouter` calls the receiver methods of the `PackageSet` to execute the command provided by the `Messages`. The results are returned to the `MessageRouter` and returned to the `HandleConnection` function through a channel that is included with every message to send the result back to the proper Go routine to be returned to the client.
