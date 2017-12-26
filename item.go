package trapi

import (
	"strings"
)

type DataType int

const (
	DATATYPE_NONE DataType = iota
	DATATYPE_STRING
	DATATYPE_NUMBER
	DATATYPE_INTEGER
	DATATYPE_BOOLEAN
	DATATYPE_OBJECT
	DATATYPE_ARRAY
	DATATYPE_BINARY
	DATATYPE_DATE
	DATATYPE_TIME
	DATATYPE_DATETIME
	DATATYPE_CUSTOM = 1000
)

func (dt DataType) String() string {
	switch dt {
	case DATATYPE_NONE:
		return "DATATYPE_NONE"
	case DATATYPE_STRING:
		return "DATATYPE_STRING"
	case DATATYPE_NUMBER:
		return "DATATYPE_NUMBER"
	case DATATYPE_INTEGER:
		return "DATATYPE_INTEGER"
	case DATATYPE_BOOLEAN:
		return "DATATYPE_BOOLEAN"
	case DATATYPE_OBJECT:
		return "DATATYPE_OBJECT"
	case DATATYPE_ARRAY:
		return "DATATYPE_ARRAY"
	case DATATYPE_BINARY:
		return "DATATYPE_BINARY"
	case DATATYPE_DATE:
		return "DATATYPE_DATE"
	case DATATYPE_TIME:
		return "DATATYPE_TIME"
	case DATATYPE_DATETIME:
		return "DATATYPE_DATETIME"
	case DATATYPE_CUSTOM:
		return "DATATYPE_CUSTOM"
	}
	return "DATATYPE_UNKNOWN"
}

type ParamType int

const (
	PARAMTYPE_UNKNOWN ParamType = iota
	PARAMTYPE_QUERY
	PARAMTYPE_URI
	PARAMTYPE_BODY
)

func (pt ParamType) String() string {
	switch pt {
	case PARAMTYPE_UNKNOWN:
		return "PARAMTYPE_UNKNOWN"
	case PARAMTYPE_QUERY:
		return "PARAMTYPE_QUERY"
	case PARAMTYPE_URI:
		return "PARAMTYPE_URI"
	case PARAMTYPE_BODY:
		return "PARAMTYPE_BODY"
	}
	return "PARAMTYPE_UNKNOWN"
}

type ResponseType int

const (
	RESPONSETYPE_UNKNOWN ResponseType = iota
	RESPONSETYPE_SUCCESS
	RESPONSETYPE_ERROR
)

func (rt ResponseType) String() string {
	switch rt {
	case RESPONSETYPE_UNKNOWN:
		return "RESPONSETYPE_UNKNOWN"
	case RESPONSETYPE_SUCCESS:
		return "RESPONSETYPE_SUCCESS"
	case RESPONSETYPE_ERROR:
		return "RESPONSETYPE_ERROR"
	}
	return "RESPONSETYPE_UNKNOWN"
}

func ParseParamType(param_type string) ParamType {
	switch param_type {
	case "query":
		return PARAMTYPE_QUERY
	case "uri":
		return PARAMTYPE_URI
	case "body":
		return PARAMTYPE_BODY
	default:
		return PARAMTYPE_UNKNOWN
	}
}

func ParseResponseType(response_type string) ResponseType {
	switch response_type {
	case "success":
		return RESPONSETYPE_SUCCESS
	case "error":
		return RESPONSETYPE_ERROR
	default:
		return RESPONSETYPE_UNKNOWN
	}
}

type ApiDataType struct {
	DataTypeName string
	DataType     DataType
	ItemType     *string
	ParentType   *string
	Description  string
	Items        map[string]*ApiDataTypeField
	ItemsOrder   []string
	Examples     []*ApiExample
	BuiltIn      bool
	Override     bool
}

func (a *ApiDataType) Clone() *ApiDataType {
	ret := &ApiDataType{
		DataTypeName: a.DataTypeName,
		DataType:     a.DataType,
		ItemType:     a.ItemType,
		ParentType:   a.ParentType,
		Description:  a.Description,
		Examples:     a.Examples,
		BuiltIn:      a.BuiltIn,
		Override:     a.Override,
	}
	if a.Items != nil {
		ret.Items = make(map[string]*ApiDataTypeField)
		for ak, av := range a.Items {
			ret.Items[ak] = &ApiDataTypeField{
				FieldName:   av.FieldName,
				Required:    av.Required,
				ApiDataType: av.ApiDataType.Clone(),
			}
		}
	}
	if a.ItemsOrder != nil {
		for _, ao := range a.ItemsOrder {
			ret.ItemsOrder = append(ret.ItemsOrder, ao)
		}
	}
	return ret
}

type ApiDataTypeField struct {
	FieldName   string
	Required    bool
	ApiDataType *ApiDataType
}

type Api struct {
	Method      string
	Path        string
	Description string
	Params      ApiParamTypeList
	Responses   []*ApiResponse
	Headers     *ApiHeaderList

	SPIB_Filename
}

type ApiList struct {
	Path     string
	SubItems []*ApiList
	Apis     []*Api
}

func (a *ApiList) Find(path string) *ApiList {
	for _, ap := range a.SubItems {
		if ap.Path == path {
			return ap
		}
	}
	return nil
}

func (a *ApiList) Add(api *Api) {

	paths := strings.Split(api.Path, "/")
	cura := a
	for _, p := range paths {
		if p != "" {
			fal := cura.Find(p)
			if fal == nil {
				fal = &ApiList{
					Path: p,
				}
				cura.SubItems = append(cura.SubItems, fal)
			}
			cura = fal
		}
	}

	cura.Apis = append(cura.Apis, api)
}

type ApiDefine struct {
	DefineType string
	Name       string
	DataType   *ApiDataType
	Examples   []*ApiExample

	SPIB_Filename
}

type ApiParam struct {
	Name     string
	DataType *ApiDataType

	SPIB_Filename
}

type ApiParamList struct {
	List  map[string]*ApiParam
	Order []string
}

type ApiParamTypeList map[ParamType]*ApiParamList

type ApiHeader struct {
	Name        string
	DataType    *ApiDataType
	Description string

	SPIB_Filename
}

type ApiHeaderList struct {
	List  map[string][]*ApiHeader
	Order []string
}

type ApiExample struct {
	ContentType string
	Description string
	Text        string

	SPIB_Filename
}

type ApiResponse struct {
	ResponseType ResponseType
	Codes        string
	ContentTypes string
	DataType     *ApiDataType

	Examples []*ApiExample
	Headers  *ApiHeaderList
}
