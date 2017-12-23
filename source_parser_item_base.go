package trapi

import (
	"strings"
)

type SPIB_Filename struct {
	Filename string
	Line     int
}

type SPIB_DataType struct {
	Name        string
	DataType    string
	Description string
	Required    bool
	Items       SPIB_DataTypeList
}

func NewSPIB_DataType(name string, datatype string, description string) SPIB_DataType {
	name = strings.TrimSpace(name)
	datatype = strings.TrimSpace(datatype)
	description = strings.TrimSpace(description)
	required := true

	if strings.HasSuffix(name, "?") {
		required = false
		name = strings.TrimSuffix(name, "?")
	}

	return SPIB_DataType{
		Name:        name,
		DataType:    datatype,
		Description: description,
		Required:    required,
	}
}

type SPIB_Text struct {
	Text string
}

type SPIB_DataTypeList []*SPIB_DataType

func (pd SPIB_DataTypeList) Find(name string) *SPIB_DataType {
	for _, p := range pd {
		if p.Name == name {
			return p
		}
	}
	return nil
}

type ISPIB_WithExamples interface {
	AppendExample(example *SourceParseItemExample)
}

type ISPIB_WithHeaders interface {
	AppendHeader(header *SourceParseItemHeader)
}
