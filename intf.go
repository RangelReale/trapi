package trapi

import (
	"io"
)

type Generator interface {
	Generate(parser *Parser, out io.Writer) error
}
