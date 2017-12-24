package trapi

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/RangelReale/gocompar"
)

var (
	ErrIgnore = errors.New("Ignore")
)

type sourceParserFile struct {
	parser   *SourceParser
	filename string
	stack    *SourceParseStack
}

func newSourceParserFile(parser *SourceParser, filename string) *sourceParserFile {
	return &sourceParserFile{
		parser:   parser,
		filename: filename,
		stack:    NewSourceParseStack(),
	}
}

var (
	reAPI = regexp.MustCompile(`@api(\w*)`)
)

func (p *sourceParserFile) parseComments(comments []*gocompar.Comment) error {

	for _, c := range comments {
		err := p.parseComment(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *sourceParserFile) finish() {

	p.stackClose()

}

func (p *sourceParserFile) stackClose() error {

	for p.stack.Len() > 0 {
		err := p.stackCloseLast()
		if err != nil {
			return err
		}
	}

	return nil

}

func (p *sourceParserFile) stackCloseUntil(itemtype SourceStackItemType) error {

	for p.stack.Len() > 0 && p.stack.Top().StackItemType != itemtype {
		err := p.stackCloseLast()
		if err != nil {
			return err
		}
	}

	return nil

}

func (p *sourceParserFile) stackCloseLast() error {

	i := p.stack.Pop()
	switch i.ItemType {
	case SPARSE_ITEM_DEFINE:
		p.parser.Defines = append(p.parser.Defines, i.Item.(*SourceParseItemDefine))
	case SPARSE_ITEM_API:
		p.parser.Apis = append(p.parser.Apis, i.Item.(*SourceParseItemApi))
	case SPARSE_ITEM_PARAM:
		api := p.stack.Top().Item.(*SourceParseItemApi)
		api.AddParam(i.Item.(*SourceParseItemParam))
	case SPARSE_ITEM_RESPONSE:
		api := p.stack.Top().Item.(*SourceParseItemApi)
		api.Responses = append(api.Responses, i.Item.(*SourceParseItemResponse))
	case SPARSE_ITEM_EXAMPLE:
		withexample := p.stack.Top().Item.(ISPIB_WithExamples)
		e := i.Item.(*SourceParseItemExample)
		if withexample == nil {
			return NewParserError(fmt.Sprintf("Top item does not support examples - cannot add example %s", e.Description), e.Filename, e.Line)
		}
		withexample.AppendExample(e)
	case SPARSE_ITEM_HEADER:
		withheader := p.stack.Top().Item.(ISPIB_WithHeaders)
		h := i.Item.(*SourceParseItemHeader)
		if withheader == nil {
			return NewParserError(fmt.Sprintf("Top item does not support headers - cannot add header %s", h.Name), h.Filename, h.Line)
		}
		withheader.AppendHeader(h)
	default:
		return NewParserError(fmt.Sprintf("Unknown item type: %v", i.Item), p.filename, 0)
	}

	return nil

}

func (p *sourceParserFile) parseComment(comment *gocompar.Comment) error {

	scan := bufio.NewScanner(strings.NewReader(comment.Text))
	line := 0
	for scan.Scan() {

		s := reAPI.FindStringSubmatch(scan.Text())

		if p.stack.Len() > 0 && p.stack.Top().StackItemType == SITEM_TEXT {
			fv := scan.Text()

			if s != nil && len(s) > 1 {
				// api tag found, end example
				err := p.stackCloseLast()
				if err != nil {
					return err
				}
			} else if fv == "" {
				// empty line found, end example
				err := p.stackCloseLast()
				if err != nil {
					return err
				}
			} else {
				pit := p.stack.Top().StackItem.(*SPIB_Text)

				if pit.Text != "" {
					pit.Text += "\n"
				}

				pit.Text += fv
			}
		}

		if s != nil && len(s) > 1 {
			//fmt.Printf("FOUND: [%s] %v\n", p.filename, s)

			var err error
			switch s[1] {
			case "Define":
				err = p.parseDefine(line, comment, scan.Text())
			case "":
				err = p.parseApi(line, comment, scan.Text())
			case "Param":
				err = p.parseParam(line, comment, scan.Text())
			case "Success", "Error", "Response":
				err = p.parseResponse(line, comment, scan.Text())
			case "Field":
				err = p.parseData(line, comment, scan.Text())
			case "Example":
				err = p.parseExample(line, comment, scan.Text())
			case "Header":
				err = p.parseHeader(line, comment, scan.Text())
			case "Ignore":
				return nil
			case "IgnoreFile":
				return ErrIgnore
			default:
				err = NewParserError(fmt.Sprintf("Unknown directive @api%s", s[1]), p.filename, comment.Line)
			}
			if err != nil {
				return err
			}
		}
		line++
	}

	if err := scan.Err(); err != nil {
		return err
	}

	return nil
}
