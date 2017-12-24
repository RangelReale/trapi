package trapi

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/RangelReale/gocompar"
)

//
// @api: Define
//

var (
	// @apiDefine (define_type) {data_type} Name Description
	reAPIDefine = regexp.MustCompile(`@apiDefine \(([^\)]+)\) \{([^}]+)\} (\S+)(.*)$`)
)

func (p *sourceParserFile) parseDefine(line int, comment *gocompar.Comment, text string) error {

	// close any pending item
	err := p.stackClose()
	if err != nil {
		return err
	}

	s := reAPIDefine.FindStringSubmatch(text)
	if s == nil || len(s) < 2 {
		return fmt.Errorf("Could not parse @apiDefine line: %s", text)
	}

	//fmt.Printf("@apiDefine: {%+v} [[[%s]]]\n", strings.Join(s[1:], ", "), text)

	newi := &SourceParseItemDefine{
		DefineType:    strings.TrimSpace(s[1]),
		SPIB_DataType: NewSPIB_DataType(s[3], s[2], s[4]),
		SPIB_Filename: SPIB_Filename{
			Filename: p.filename,
			Line:     comment.Line + line,
		},
	}

	p.stack.Push(&SourceStackData{
		ItemType:      SPARSE_ITEM_DEFINE,
		Item:          newi,
		StackItemType: SITEM_DATATYPE,
		StackItem:     &newi.SPIB_DataType,
	})

	return nil
}

//
// @api: API
//

var (
	// @api {method} path Description
	reAPIAPI       = regexp.MustCompile(`@api \{([^}]+)\} (\S+)(.*)$`)
	reAPIAPIParams = regexp.MustCompile(`<([^>]+)>`)
)

func (p *sourceParserFile) parseApi(line int, comment *gocompar.Comment, text string) error {

	// close any pending item
	err := p.stackClose()
	if err != nil {
		return err
	}

	s := reAPIAPI.FindStringSubmatch(text)
	if s == nil || len(s) < 2 {
		return fmt.Errorf("Could not parse @api line: %s", text)
	}

	//fmt.Printf("@api: {%+v} [[[%s]]]\n", strings.Join(s[1:], ", "), text)

	newi := &SourceParseItemApi{
		Method:      strings.TrimSpace(s[1]),
		Path:        strings.TrimSpace(s[2]),
		Description: strings.TrimSpace(s[3]),
		SPIB_Filename: SPIB_Filename{
			Filename: p.filename,
			Line:     comment.Line + line,
		},
	}

	// extract params from path
	params := reAPIAPIParams.FindAllStringSubmatch(newi.Path, -1)
	for _, pi := range params {
		newi.Params = append(newi.Params, &SourceParseItemParam{
			ParamType:     "uri",
			Name:          pi[1],
			SPIB_DataType: NewSPIB_DataType("param", "String", ""),
			SPIB_Filename: SPIB_Filename{
				Filename: p.filename,
				Line:     comment.Line + line,
			},
		})
	}

	p.stack.Push(&SourceStackData{
		ItemType:      SPARSE_ITEM_API,
		Item:          newi,
		StackItemType: SITEM_API,
		StackItem:     newi,
	})

	return nil
}

//
// @api: Param
//

var (
	// @apiParam param_type {data_type} Name Description
	reAPIParam = regexp.MustCompile(`@apiParam (\S+) \{([^}]+)\} (\S+)(.*)$`)
)

func (p *sourceParserFile) parseParam(line int, comment *gocompar.Comment, text string) error {

	// must have an "Api" item at top
	err := p.stackCloseUntil(SITEM_API)
	if err != nil {
		return err
	}

	if p.stack.Top() == nil || p.stack.Top().ItemType != SPARSE_ITEM_API {
		return fmt.Errorf("@apiParam must come after an @api: %s", text)
	}

	s := reAPIParam.FindStringSubmatch(text)
	if s == nil || len(s) < 2 {
		return fmt.Errorf("Could not parse @apiParam line: %s", text)
	}

	//fmt.Printf("@apiParam: {%+v} [[[%s]]]\n", strings.Join(s[1:], ", "), text)

	item_param := &SourceParseItemParam{
		ParamType:     strings.TrimSpace(s[1]),
		Name:          strings.TrimSuffix(strings.TrimSpace(s[3]), "?"),
		SPIB_DataType: NewSPIB_DataType("param", strings.TrimSpace(s[2]), strings.TrimSpace(s[4])),
		SPIB_Filename: SPIB_Filename{
			Filename: p.filename,
			Line:     comment.Line + line,
		},
	}

	p.stack.Push(&SourceStackData{
		ItemType:      SPARSE_ITEM_PARAM,
		Item:          item_param,
		StackItemType: SITEM_DATATYPE,
		StackItem:     &item_param.SPIB_DataType,
	})

	return nil
}

//
// @api: Response
//

var (
	// @api{Success,Error} codes content_types {data_type} Description
	reAPIResponse = regexp.MustCompile(`@api(\w+) (\S+) (\S+) \{([^}]+)\} (.*)$`)
)

func (p *sourceParserFile) parseResponse(line int, comment *gocompar.Comment, text string) error {

	// must have an "Api" item at top
	err := p.stackCloseUntil(SITEM_API)
	if err != nil {
		return err
	}

	if p.stack.Top() == nil || p.stack.Top().ItemType != SPARSE_ITEM_API {
		return fmt.Errorf("@apiSuccess/@apiError must come after an @api: %s", text)
	}

	s := reAPIResponse.FindStringSubmatch(text)
	if s == nil || len(s) < 2 {
		return fmt.Errorf("Could not parse @apiSuccess/@apiError line: %s", text)
	}

	//fmt.Printf("@apiResponse: {%+v} [[[%s]]]\n", strings.Join(s[1:], ", "), text)

	newi := &SourceParseItemResponse{
		ResponseType:  strings.ToLower(s[1]),
		Codes:         strings.TrimSpace(s[2]),
		ContentTypes:  strings.TrimSpace(s[3]),
		SPIB_DataType: NewSPIB_DataType("response", strings.TrimSpace(s[4]), strings.TrimSpace(s[5])),
		SPIB_Filename: SPIB_Filename{
			Filename: p.filename,
			Line:     comment.Line + line,
		},
	}

	p.stack.Push(&SourceStackData{
		ItemType:      SPARSE_ITEM_RESPONSE,
		Item:          newi,
		StackItemType: SITEM_DATATYPE,
		StackItem:     &newi.SPIB_DataType,
	})

	return nil
}

//
// @api: Data
//

var (
	// @apiField {data_type} Name Description
	reAPIField = regexp.MustCompile(`@apiField \{([^}]+)\} (\S+)(.*)$`)
)

func (p *sourceParserFile) parseData(line int, comment *gocompar.Comment, text string) error {

	// must have an "DataType" item at top
	if p.stack.Top() == nil || p.stack.Top().StackItemType != SITEM_DATATYPE {
		return fmt.Errorf("@apiField must come after an datatype definition: %s", text)
	}

	s := reAPIField.FindStringSubmatch(text)
	if s == nil || len(s) < 2 {
		return fmt.Errorf("Could not parse @apiField line: %s", text)
	}

	//fmt.Printf("@apiField: {%+v} [[[%s]]]\n", strings.Join(s[1:], ", "), text)

	datatype := p.stack.Top().StackItem.(*SPIB_DataType)

	sub := strings.Split(s[2], ".")
	curdt := datatype
	for subct, subname := range sub {
		if curdt.Items == nil {
			curdt.Items = make(SPIB_DataTypeList, 0)
		}
		newi := curdt.Items.Find(subname)
		if newi == nil {
			datatype := strings.TrimSpace(s[1])
			description := strings.TrimSpace(s[3])
			if subct < len(sub)-1 {
				datatype = "Object"
				description = ""
			}

			xnewi := NewSPIB_DataType(subname, datatype, description)
			newi = &xnewi
			curdt.Items = append(curdt.Items, newi)
		}
		curdt = newi
	}

	return nil
}

//
// @api: Example
//

var (
	// @apiExample {content_type} Description
	reAPIExample = regexp.MustCompile(`@apiExample \{([^}]+)\} (.*)$`)
)

func (p *sourceParserFile) parseExample(line int, comment *gocompar.Comment, text string) error {

	/*
		// must have an "Response" item at top
		if p.stack.Top() == nil || p.stack.Top().Item != ITEM_RESPONSE {
			return fmt.Errorf("@apiExample must come after an @apiSuccess/@apiError: %s", text)
		}
	*/

	s := reAPIExample.FindStringSubmatch(text)
	if s == nil || len(s) < 2 {
		return fmt.Errorf("Could not parse @apiExample line: %s", text)
	}

	//fmt.Printf("@apiExample: {%+v} [[[%s]]]\n", strings.Join(s[1:], ", "), text)

	newi := &SourceParseItemExample{
		ContentType: strings.TrimSpace(s[1]),
		Description: strings.TrimSpace(s[2]),
		SPIB_Text: SPIB_Text{
			Text: "",
		},
		SPIB_Filename: SPIB_Filename{
			Filename: p.filename,
			Line:     comment.Line + line,
		},
	}

	p.stack.Push(&SourceStackData{
		ItemType:      SPARSE_ITEM_EXAMPLE,
		Item:          newi,
		StackItemType: SITEM_TEXT,
		StackItem:     &newi.SPIB_Text,
	})

	return nil
}

//
// @api: Header
//

var (
	// @apiHeader {data_type} Name Description
	reAPIHeader = regexp.MustCompile(`@apiHeader \{([^}]+)\} (\S+)(.*)$`)
)

func (p *sourceParserFile) parseHeader(line int, comment *gocompar.Comment, text string) error {

	s := reAPIHeader.FindStringSubmatch(text)
	if s == nil || len(s) < 2 {
		return fmt.Errorf("Could not parse @apiHeader line: %s", text)
	}

	//fmt.Printf("@apiHeader: {%+v} [[[%s]]]\n", strings.Join(s[1:], ", "), text)

	newi := &SourceParseItemHeader{
		DataType:    strings.TrimSpace(s[1]),
		Name:        strings.TrimSpace(s[2]),
		Description: strings.TrimSpace(s[3]),
		SPIB_Filename: SPIB_Filename{
			Filename: p.filename,
			Line:     comment.Line + line,
		},
	}

	p.stack.Push(&SourceStackData{
		ItemType:      SPARSE_ITEM_HEADER,
		Item:          newi,
		StackItemType: SITEM_NONE,
	})

	return nil
}
