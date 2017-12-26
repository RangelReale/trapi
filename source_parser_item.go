package trapi

type SourceParseItemType int

const (
	SPARSE_ITEM_DEFINE SourceParseItemType = iota
	SPARSE_ITEM_API
	SPARSE_ITEM_PARAM
	SPARSE_ITEM_RESPONSE
	SPARSE_ITEM_EXAMPLE
	SPARSE_ITEM_HEADER
)

func (sp SourceParseItemType) String() string {
	switch sp {
	case SPARSE_ITEM_DEFINE:
		return "SPARSE_ITEM_DEFINE"
	case SPARSE_ITEM_API:
		return "SPARSE_ITEM_API"
	case SPARSE_ITEM_PARAM:
		return "SPARSE_ITEM_PARAM"
	case SPARSE_ITEM_RESPONSE:
		return "SPARSE_ITEM_RESPONSE"
	case SPARSE_ITEM_EXAMPLE:
		return "SPARSE_ITEM_EXAMPLE"
	case SPARSE_ITEM_HEADER:
		return "SPARSE_ITEM_HEADER"
	}
	return "SPARSE_ITEM_UNKNOWN"
}

//
// Source Parse Item: API
//
type SourceParseItemApi struct {
	SPIB_Filename

	Method      string
	Path        string
	Description string

	Params    []*SourceParseItemParam
	Headers   []*SourceParseItemHeader
	Responses []*SourceParseItemResponse
}

func (a *SourceParseItemApi) AddParam(param *SourceParseItemParam) {
	if a.Params != nil {
		// if uri param exists, replace
		if param.ParamType == "uri" {
			for _, p := range a.Params {
				if p.ParamType == param.ParamType && p.Name == param.Name {
					*p = *param
					return
				}
			}
		}
	}
	a.Params = append(a.Params, param)
}

func (a *SourceParseItemApi) AppendHeader(header *SourceParseItemHeader) {
	a.Headers = append(a.Headers, header)
}

//
// Source Parse Item: RESPONSE
//
type SourceParseItemResponse struct {
	SPIB_Filename
	SPIB_DataType

	ResponseType string
	Codes        string
	ContentTypes string

	Headers  []*SourceParseItemHeader
	Examples []*SourceParseItemExample
}

func (pr *SourceParseItemResponse) AppendExample(example *SourceParseItemExample) {
	pr.Examples = append(pr.Examples, example)
}

func (pr *SourceParseItemResponse) AppendHeader(header *SourceParseItemHeader) {
	pr.Headers = append(pr.Headers, header)
}

//
// Source Parse Item: Param
//
type SourceParseItemParam struct {
	SPIB_Filename
	SPIB_DataType

	ParamType string
	Name      string

	Examples []*SourceParseItemExample
}

func (pr *SourceParseItemParam) AppendExample(example *SourceParseItemExample) {
	pr.Examples = append(pr.Examples, example)
}

//
// Source Parse Item: DEFINE
//
type SourceParseItemDefine struct {
	SPIB_Filename
	SPIB_DataType

	DefineType string
	Examples   []*SourceParseItemExample
}

func (pr *SourceParseItemDefine) AppendExample(example *SourceParseItemExample) {
	pr.Examples = append(pr.Examples, example)
}

//
// Source Parse Item: HEADER
//
type SourceParseItemHeader struct {
	SPIB_Filename

	DataType    string
	Name        string
	Description string
}

//
// SourceParse Item: EXAMPLE
//
type SourceParseItemExample struct {
	SPIB_Filename
	SPIB_Text

	ContentType string
	Description string
}
