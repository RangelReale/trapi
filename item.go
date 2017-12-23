package trapi

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
	DATATYPE_CUSTOM = 1000
)

type ParamType int

const (
	PARAMTYPE_UNKNOWN ParamType = iota
	PARAMTYPE_QUERY
	PARAMTYPE_URI
	PARAMTYPE_BODY
)

type ResponseType int

const (
	RESPONSETYPE_UNKNOWN ResponseType = iota
	RESPONSETYPE_SUCCESS
	RESPONSETYPE_ERROR
)

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
	Name         string
	DataType     DataType
	ItemType     *ApiDataType
	OriginalType string
	Description  string
	Required     bool
	Items        map[string]*ApiDataType
	ItemsOrder   []string
	Examples     []*ApiExample
}

func (a *ApiDataType) Clone() *ApiDataType {
	ret := &ApiDataType{
		Name:         a.Name,
		DataType:     a.DataType,
		ItemType:     a.ItemType,
		OriginalType: a.OriginalType,
		Description:  a.Description,
		Required:     a.Required,
		Examples:     a.Examples,
	}
	if a.Items != nil {
		ret.Items = make(map[string]*ApiDataType)
		for ak, av := range a.Items {
			ret.Items[ak] = av.Clone()
		}
	}
	if a.ItemsOrder != nil {
		for _, ao := range a.ItemsOrder {
			ret.ItemsOrder = append(ret.ItemsOrder, ao)
		}
	}
	return ret
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