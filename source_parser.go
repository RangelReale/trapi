package trapi

import (
	"github.com/RangelReale/gocompar"
)

type SourceParser struct {
	gcp *gocompar.Parser
	tags []string

	Defines []*SourceParseItemDefine
	Apis    []*SourceParseItemApi
}

func NewSourceParser(gcp *gocompar.Parser) *SourceParser {
	return &SourceParser{
		gcp: gcp,
	}
}

func (p *SourceParser) AddTag(tag string) {
	p.tags = append(p.tags, tag)
}

func (p *SourceParser) AddTags(tags []string) {
	p.tags = append(p.tags, tags...)
}

func (p *SourceParser) Process() error {

	for _, f := range p.gcp.Comments {

		fp := newSourceParserFile(p, f.Filename)
		err := fp.parseComments(f.Comments)
		if err != nil {
			if err != ErrIgnore {
				return err
			}
		}

		fp.finish()

	}

	return nil
}
