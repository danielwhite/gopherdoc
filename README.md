# Gopherdoc

Gopherdoc runs a Gopher server providing documentation for Go programs.

At this stage, it runs documentation for all packages contained with
the `$GOROOT` and `$GODOC` paths. It offers no constraints on inbound
requests, so it may be unwise to run on the public internet.

Documentation is at: at: https://godoc.org/whitehouse.id.au/gopherdoc.

## Usage

To run the Gopher server accessible at the host `go.example.com`:

```
$ go get whitehouse.id.au/gopherdoc
$ gopherdoc -host go.example.com
```
