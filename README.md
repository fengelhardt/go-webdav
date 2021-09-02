# A CLI WebDAV File Server Written in go

WebDAV is a common file server protocol supported by many file browers.

This lean but efficient server runs from the command line and holds all
its data just within a folder in your directory tree.
No database is required.

Features are simple but powerful:

- TLS support (https)
- SHA3 hashes used in its password file
- Multi-user functionality
- Configuration via simple config files
- All data is stored in the file system
- User credentials can be generated with the `go-webdav-hash` tool
- File locking (not persistent between individual server invocations)

## Running the Test Configuration

Clone this repo and build the code.

	go get -u ./...
	go build .

Create the folder where the server should store all its data.
(Temporary in tis case.)

	mkdir /tmp/srv

Now run the server by telling it where the config file is:

	./go-webdav -c test/go-webdav.hjson

You can now browse `webdavs://localhost:9999/webdav` with your avorite WebDAV client.

You have to accept the test license, which is in no way secure, but there for testing purposes.
Use either "alice" (pwd: alice) or "bob" (pwd: bob) for credentials.

You can see now that in the `/tmp/srv` folder the uploaded files are placed,
with a separate subdirectory for each user.

## Configuration

There is a general config file and a separate password file for user administration.
Both use a human readable json dialect.

For documentation, see the examples in the `test/` folder.

To generate the password hashes, use the `go-webdav-hash` tool.

## TODO / What is Missing

- Quota implementation
- Separate access and error logs
- Logging via syslog
- More sanity checks for the config files
- HTTPS redirection. Currently only HTTPS is supported, no HTTP.

## Tested Clients

- dolphin
