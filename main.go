package main

import (
	"flag"
	"go/build"
	"log"
	"net"
	"net/textproto"
	"path/filepath"
	"runtime"
	"text/template"
	"time"

	"golang.org/x/tools/godoc"
	"golang.org/x/tools/godoc/static"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

// Maximum amount of time to wait when reading from a client before
// timing out.
const readTimeout = time.Second * 10

var (
	hostFlag    = flag.String("host", "localhost", "gopher service hostname")
	addrFlag    = flag.String("addr", ":gopher", "gopher service address (e.g. ':70')")
	gorootFlag  = flag.String("goroot", runtime.GOROOT(), "Go root directory")
	verboseFlag = flag.Bool("v", false, "verbose mode")
)

func main() {
	flag.Parse()

	// TODO: Limit concurrent access to filesystem.
	fs := vfs.NameSpace{}
	rootfs := vfs.OS(*gorootFlag)
	fs.Bind("/", rootfs, "/", vfs.BindReplace)
	fs.Bind("/lib/godoc", mapfs.New(static.Files), "/", vfs.BindReplace)

	// Bind $GOPATH trees into Go root.
	for _, p := range filepath.SplitList(build.Default.GOPATH) {
		fs.Bind("/src", vfs.OS(p), "/src", vfs.BindAfter)
	}

	corpus := godoc.NewCorpus(fs)
	corpus.Verbose = *verboseFlag
	if err := corpus.Init(); err != nil {
		log.Fatal(err)
	}
	pres := godoc.NewPresentation(corpus)
	pres.PackageText = readTemplate(fs, pres, "package.txt")

	ln, err := net.Listen("tcp", *addrFlag)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("listening on %s", ln.Addr())

	handler := &handler{
		pres: pres,
		host: *hostFlag,
		port: getPort(ln.Addr()),
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}
		go serve(conn, handler)
	}
}

// get the port as a number from an address; this allows for callers
// to use named ports in the address.
func getPort(addr net.Addr) int {
	checkErr := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	_, portStr, err := net.SplitHostPort(addr.String())
	checkErr(err)

	port, err := net.LookupPort(addr.Network(), portStr)
	checkErr(err)

	return port
}

func readTemplate(fs vfs.FileSystem, pres *godoc.Presentation, name string) *template.Template {
	path := "lib/godoc/" + name

	// template package can't read directly from a vfs, so we need
	// to read the data ourselves.
	data, err := vfs.ReadFile(fs, path)
	if err != nil {
		log.Fatal("readTemplate: ", err)
	}

	return template.Must(template.New(name).Funcs(pres.FuncMap()).Parse(string(data)))
}

// Serve a new connection.
func serve(conn net.Conn, handler *handler) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("panic serving %v: %v\n%s", conn.RemoteAddr(), err, buf)
		}
		conn.Close()
	}()

	// Don't let clients hog connections forever.
	conn.SetReadDeadline(time.Now().Add(readTimeout))

	handler.Handle(textproto.NewConn(conn))
}
