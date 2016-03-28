# do-code-challenge
## Design Document

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
