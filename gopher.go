package main

import (
	"fmt"
	"io"
)

type ItemType byte

const (
	TextfileItem ItemType = '0'
	MenuItem              = '1'
	InfoItem              = 'i'
)

type DirEntity struct {
	Type     ItemType
	Name     string
	Selector string
	Host     string
	Port     int
}

func (e DirEntity) String() string {
	return fmt.Sprintf("%c%s\t%s\t%s\t%d", e.Type, e.Name, e.Selector, e.Host, e.Port)
}

type MenuEntity []DirEntity

func (e MenuEntity) Write(w io.Writer) {
	for _, entity := range e {
		fmt.Fprintf(w, "%s\n", entity)
	}
}
