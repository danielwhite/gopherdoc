/*
Gopherdoc runs a Gopher server providing documentation for Go programs.

Documentation is served for all packages contained with the $GOROOT
and $GODOC paths.

The addr flag is used to control the interface and port binding.

	gopherdoc -addr 0.0.0.0:gopher

The host flag determines the host used in directory entities. If this
server were to be accessible at go.example.com, then the following
would ensure the server generates menu entries that are accessible:

	gopherdoc -host go.example.com
*/
package main // import "whitehouse.id.au/gopherdoc"
