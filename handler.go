package main

import (
	"fmt"
	"io"
	"log"
	"net/textproto"
	pathpkg "path"
	"strings"

	"golang.org/x/tools/godoc"
)

type handler struct {
	pres *godoc.Presentation
	host string // public hostname of server
	port int    // public port of server
}

// Handle a client request.
func (h *handler) Handle(conn *textproto.Conn) {
	// Read the line for the transaction; this is typically a selector.
	path, err := conn.ReadLine()
	if err != nil {
		return
	}

	w := conn.DotWriter()
	defer w.Close()

	// Write a documentation page if selector is a reference to one.
	if strings.HasPrefix(path, "doc:") {
		h.writeDoc(w, h.getPageInfo(path[4:]))
		return
	}

	// Otherwise build up a menu entity for the response.
	info := h.getPageInfo(path)
	h.makeMenu(info, path).Write(w)
}

// getPageInfo reads information for displaying a page of information
// about the package referenced by path.
func (h *handler) getPageInfo(path string) *godoc.PageInfo {
	const pageInfoMode = godoc.NoFiltering | godoc.NoHTML
	info := h.pres.GetPkgPageInfo("/src/"+path, path, pageInfoMode)
	if info.Err != nil {
		log.Panicf("getPageInfo: %s", info.Err)
	}
	log.Printf("directory: %s", info.Dirname)
	return info
}

// write the textual form of package documentation for info to w.
func (h *handler) writeDoc(w io.Writer, info *godoc.PageInfo) error {
	return h.pres.PackageText.Execute(w, info)
}

// make a gopher menu entity for the given page of info.
func (h *handler) makeMenu(info *godoc.PageInfo, parentPath string) (entities MenuEntity) {
	// add a textfile entry for documentation if available
	if info.PDoc != nil {
		entities = append(entities, DirEntity{
			Type:     TextfileItem,
			Name:     "Package Documentation",
			Selector: "doc:" + parentPath,
			Host:     h.host,
			Port:     h.port,
		})
	}

	// add a directory entry for all immediate sub-packages
	if info.Dirs == nil {
		return
	}

	for _, dir := range info.Dirs.List {
		if dir.Depth > 0 {
			continue
		}
		entities = append(entities, h.makeDirEntity(&dir, parentPath))
	}
	return
}

// make a directory entity for based on the sub-package entry.
func (h *handler) makeDirEntity(dir *godoc.DirEntry, parentPath string) DirEntity {
	path := pathpkg.Join(parentPath, dir.Path)

	entity := DirEntity{
		Type:     MenuItem,
		Name:     pathpkg.Base(path),
		Selector: path,
		Host:     h.host,
		Port:     h.port,
	}

	if dir.Synopsis != "" {
		entity.Name = fmt.Sprintf("%-11s: %s", entity.Name, dir.Synopsis)
	}

	return entity
}
