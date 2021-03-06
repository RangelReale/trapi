package trapi

import (
	"fmt"
	"strings"

	"github.com/RangelReale/gocompar"
)

type Parser struct {
	files      []string
	dirs       []string
	gcp        *gocompar.Parser
	apidefload []*SourceParseItemDefine
	tags []string

	DataTypes  map[string]*ApiDataType
	ApiDefines []*ApiDefine
	Apis       []*Api
}

func NewParser(gcp *gocompar.Parser) *Parser {
	ret := &Parser{
		gcp: gcp,
		DataTypes: map[string]*ApiDataType{
			"String": &ApiDataType{
				DataTypeName: "String",
				DataType:     DATATYPE_STRING,
				BuiltIn:      true,
			},
			"Number": &ApiDataType{
				DataTypeName: "Number",
				DataType:     DATATYPE_NUMBER,
				BuiltIn:      true,
			},
			"Integer": &ApiDataType{
				DataTypeName: "Integer",
				DataType:     DATATYPE_INTEGER,
				BuiltIn:      true,
			},
			"Boolean": &ApiDataType{
				DataTypeName: "Boolean",
				DataType:     DATATYPE_BOOLEAN,
				BuiltIn:      true,
			},
			"Date": &ApiDataType{
				DataTypeName: "Date",
				DataType:     DATATYPE_DATE,
				BuiltIn:      true,
			},
			"Time": &ApiDataType{
				DataTypeName: "Time",
				DataType:     DATATYPE_TIME,
				BuiltIn:      true,
			},
			"DateTime": &ApiDataType{
				DataTypeName: "DateTime",
				DataType:     DATATYPE_DATETIME,
				BuiltIn:      true,
			},
			"Object": &ApiDataType{
				DataTypeName: "Object",
				DataType:     DATATYPE_OBJECT,
				BuiltIn:      true,
			},
			"Array": &ApiDataType{
				DataTypeName: "Array",
				DataType:     DATATYPE_ARRAY,
				BuiltIn:      true,
			},
			"Binary": &ApiDataType{
				DataTypeName: "Binary",
				DataType:     DATATYPE_BINARY,
				BuiltIn:      true,
			},
		},
	}

	return ret
}

func (p *Parser) AddFile(filename string) {
	p.files = append(p.files, filename)
}

func (p *Parser) AddDir(dir string) {
	p.dirs = append(p.dirs, dir)
}

func (p *Parser) AddTag(tag string) {
	p.tags = append(p.tags, tag)
}

func (p *Parser) AddTags(tags []string) {
	p.tags = append(p.tags, tags...)
}

func (p *Parser) Parse() error {

	var err error
	for _, f := range p.files {
		err = p.gcp.ParseFile(f)
		if err != nil {
			return err
		}
	}
	for _, f := range p.dirs {
		err = p.gcp.ParseDir(f)
		if err != nil {
			return err
		}
	}

	sp := NewSourceParser(p.gcp)
	sp.AddTags(p.tags)
	err = sp.Process()
	if err != nil {
		return err
	}
	return p.ParseSource(sp)
}

func (p *Parser) ParseSource(sp *SourceParser) error {

	// Do multiple passes to load all dependent types
	for dct := 0; ; dct++ {
		ctconv, ctmiss, err := p.parseSourceDefinesPass(sp)
		if err != nil {
			return err
		}

		//fmt.Printf("*** PASS %d [ctconv:%d ctmiss:%d]\n", dct, ctconv, ctmiss)

		if ctmiss == 0 {
			break
		}

		if ctconv == 0 {
			miss_def := make([]string, 0)
			for _, d := range sp.Defines {
				if _, founddt := p.DataTypes[d.Name]; !founddt {
					miss_def = append(miss_def, d.Name)
				}
			}

			return fmt.Errorf("Could not resolve all api references: missing [%s]", strings.Join(miss_def, ","))
		}
	}

	// load defines
	for _, srcdefine := range p.apidefload {

		newi := &ApiDefine{
			DefineType:    srcdefine.DefineType,
			Name:          srcdefine.Name,
			SPIB_Filename: srcdefine.SPIB_Filename,
		}

		// parse data type
		dt, ctmiss, err := p.parseSourceDataType(&srcdefine.SPIB_DataType, nil, false, false)
		if err != nil {
			return err
		}
		if dt == nil || ctmiss > 0 {
			return NewParserError(fmt.Sprintf("Unknown param datatype %s", srcdefine.DataType), srcdefine.Filename, srcdefine.Line)
		}

		newi.DataType = dt

		p.ApiDefines = append(p.ApiDefines, newi)

	}

	// load apis
	for _, srcapi := range sp.Apis {

		newi := &Api{
			Method:        srcapi.Method,
			Path:          srcapi.Path,
			Description:   srcapi.Description,
			SPIB_Filename: srcapi.SPIB_Filename,
		}

		//
		// Params
		//
		for _, srcapiparam := range srcapi.Params {

			// parse param type
			pt := ParseParamType(srcapiparam.ParamType)
			if pt == PARAMTYPE_UNKNOWN {
				return NewParserError(fmt.Sprintf("Unknown param type %s", srcapiparam.ParamType), srcapiparam.Filename, srcapiparam.Line)
			}

			// parse data type
			dt, ctmiss, err := p.parseSourceDataType(&srcapiparam.SPIB_DataType, nil, false, false)
			if err != nil {
				return err
			}
			if dt == nil || ctmiss > 0 {
				return NewParserError(fmt.Sprintf("Unknown param datatype %s", srcapiparam.DataType), srcapiparam.Filename, srcapiparam.Line)
			}

			type _dtitem struct {
				Name     string
				DataType *ApiDataType
				Examples []*ApiExample
			}
			dtlist := make([]*_dtitem, 0)
			if pt == PARAMTYPE_QUERY && dt.DataType == DATATYPE_OBJECT {
				// expand keys into parameters
				for _, ppi := range dt.ItemsOrder {
					if dt.Items[ppi].ApiDataType.DataType == DATATYPE_OBJECT {
						return NewParserError(fmt.Sprintf("Only one level of indirection is supported in query param %s", srcapiparam.Name), srcapiparam.Filename, srcapiparam.Line)
					}

					dtlist = append(dtlist, &_dtitem{dt.Items[ppi].FieldName, dt.Items[ppi].ApiDataType, nil})
				}
			} else {

				sae := &_dtitem{
					Name:     srcapiparam.Name,
					DataType: dt,
				}
				if len(srcapiparam.Examples) > 0 {
					p.parseApiExampleList(srcapiparam.Examples, &sae.Examples)
				}
				dtlist = append(dtlist, sae)
			}

			for _, procdt := range dtlist {

				newip := &ApiParam{
					Name:          procdt.Name,
					DataType:      procdt.DataType,
					Examples:      procdt.Examples,
					SPIB_Filename: srcapiparam.SPIB_Filename,
				}

				if newi.Params == nil {
					newi.Params = make(ApiParamTypeList)
				}
				if _, ptok := newi.Params[pt]; !ptok {
					newi.Params[pt] = &ApiParamList{
						List: make(map[string]*ApiParam),
					}
				}

				// check if already exists
				if _, pexists := newi.Params[pt].List[newip.Name]; pexists {
					return NewParserError(fmt.Sprintf("Param '%s' (%s) already exists in api '%s'", newip.Name, srcapiparam.ParamType, srcapi.Path), srcapiparam.Filename, srcapiparam.Line)
				}

				newi.Params[pt].List[newip.Name] = newip
				newi.Params[pt].Order = append(newi.Params[pt].Order, newip.Name)
			}
		}

		//
		// Headers
		//
		if len(srcapi.Headers) > 0 {

			newi.Headers = &ApiHeaderList{
				List: make(map[string][]*ApiHeader),
			}

			err := p.parseApiHeaderList(srcapi.Headers, newi.Headers)
			if err != nil {
				return err
			}
		}

		//
		// Responses
		//
		for _, srcapiresp := range srcapi.Responses {

			// parse response type
			rt := ParseResponseType(srcapiresp.ResponseType)

			// parse data type
			dt, ctmiss, err := p.parseSourceDataType(&srcapiresp.SPIB_DataType, nil, false, false)
			if err != nil {
				return NewParserError(fmt.Sprintf("Error parsing response datatype %s [%s]", srcapiresp.DataType, err.Error()), srcapiresp.Filename, srcapiresp.Line)
			}
			if dt == nil || ctmiss > 0 {
				return NewParserError(fmt.Sprintf("Unknown response datatype %s", srcapiresp.DataType), srcapiresp.Filename, srcapiresp.Line)
			}

			codes := strings.Split(srcapiresp.Codes, ",")
			contenttypes := strings.Split(srcapiresp.ContentTypes, ",")
			if len(contenttypes) == 0 {
				contenttypes = append(contenttypes, "-")
			}

			for _, c_code := range codes {
				for _, c_contenttype := range contenttypes {

					newir := &ApiResponse{
						ResponseType: rt,
						DataType:     dt,
					}

					// response headers
					if len(srcapiresp.Headers) > 0 {
						newir.Headers = &ApiHeaderList{
							List: make(map[string][]*ApiHeader),
						}

						err := p.parseApiHeaderList(srcapiresp.Headers, newir.Headers)
						if err != nil {
							return err
						}
					}

					// response examples
					if len(srcapiresp.Examples) > 0 {
						p.parseApiExampleList(srcapiresp.Examples, &newir.Examples)
					}

					if newi.Responses == nil {
						newi.Responses = &ApiResponseList{
							List: make(map[string][]*ApiResponseBody),
						}
					}
					if _, fnd := newi.Responses.List[c_code]; !fnd {
						newi.Responses.List[c_code] = make([]*ApiResponseBody, 0)
					}

					newi.Responses.List[c_code] = append(newi.Responses.List[c_code], &ApiResponseBody{
						ContentType: c_contenttype,
						ApiResponse: newir,
					})
				}
			}

		}

		p.Apis = append(p.Apis, newi)

	}

	return nil
}

func (p *Parser) BuildApiList() *ApiList {
	ret := &ApiList{
		Path: "/",
	}
	for _, al := range p.Apis {
		ret.Add(al)
	}
	return ret
}

func (p *Parser) parseSourceDefinesPass(sp *SourceParser) (ctconv int, ctmiss int, err error) {

	ctconv = 0
	ctmiss = 0
	err = nil

	curdefined := make(map[string]bool)

	for _, d := range sp.Defines {

		if _, founddt := p.DataTypes[d.Name]; founddt {
			if _, curdef := curdefined[d.Name]; curdef {
				return 0, 0, NewParserError(fmt.Sprintf("Datatype %s was already defined", d.Name), d.Filename, d.Line)
			}
			continue
		}

		// add a temporary datatype to allow recursivity
		p.DataTypes[d.Name] = &ApiDataType{
			DataTypeName: d.Name,
			DataType:     DATATYPE_NONE,
			BuiltIn:      false,
			Override:     false,
		}

		dt, pctmiss, err := p.parseSourceDataType(&d.SPIB_DataType, nil, true, true)
		// delete temporary
		delete(p.DataTypes, d.Name)
		if err != nil {
			return 0, 0, err
		}

		ctmiss += pctmiss
		// if not found or missed some data type, leave for next pass
		if dt != nil && pctmiss == 0 {
			ctconv++
			p.DataTypes[d.Name] = dt.Clone()
			// build a list to order by load dependency
			p.apidefload = append(p.apidefload, d)

			// response examples
			if len(d.Examples) > 0 {
				p.parseApiExampleList(d.Examples, &p.DataTypes[d.Name].Examples)
			}

			curdefined[d.Name] = true
		}

	}

	return
}

func (p *Parser) getDataType(datatype string) (*ApiDataType, error) {
	if dt, ok := p.DataTypes[datatype]; ok {
		return dt, nil
	}
	return nil, fmt.Errorf("Unknown datatype '%s'", datatype)
}

func (p *Parser) parseSourceDataType(b *SPIB_DataType, rootb *SPIB_DataType, is_define bool, is_checkpass bool) (adt *ApiDataType, ctmiss int, err error) {
	if rootb == nil {
		rootb = b
	}

	is_array := false
	sdt := b.DataType
	if strings.HasSuffix(sdt, "[]") {
		is_array = true
		sdt = strings.TrimSuffix(sdt, "[]")
	}

	dt, ok := p.DataTypes[sdt]
	if ok && is_array {
		array_parenttype := "Array"
		a_itemtype := &sdt
		if dt.DataType == DATATYPE_OBJECT && dt.BuiltIn {
			a_itemtype = nil
		}
		return &ApiDataType{
			//Name:         b.Name,
			DataType:    DATATYPE_ARRAY,
			ItemType:    a_itemtype,
			ParentType:  &array_parenttype,
			Description: b.Description,
			Override:    true,
		}, 0, nil
	}

	if ok && dt.DataType != DATATYPE_OBJECT {
		ret := dt.Clone()
		ret.Description = b.Description
		if is_define {
			ret.DataTypeName = b.Name
			ret.Override = true
			ret.BuiltIn = false
		}
		return ret, 0, nil
	}

	if ok && dt.DataType == DATATYPE_OBJECT {

		ret := dt.Clone()
		ret.Description = b.Description

		if !is_define && (b.Items == nil || len(b.Items) == 0) {
			return ret, 0, nil
		}

		if is_define {
			ret.DataTypeName = b.Name
			ret.BuiltIn = false
		}
		ret.ParentType = &sdt
		ret.Override = true
		ret.OverrideItems = nil

		if b.Items != nil && len(b.Items) > 0 {
			if ret.Items == nil {
				ret.Items = make(map[string]*ApiDataTypeField)
			}
			for _, it := range b.Items {

				_, foundi := ret.Items[it.Name]

				newit, newctmiss, err := p.parseSourceDataType(it, rootb, false, is_checkpass)
				if err != nil {
					return nil, newctmiss, err
				}
				ctmiss += newctmiss
				if newctmiss == 0 {
					newifield := &ApiDataTypeField{
						FieldName:   it.Name,
						Required:    it.Required,
						ApiDataType: newit,
					}
					ret.Items[it.Name] = newifield
					if !foundi {
						ret.ItemsOrder = append(ret.ItemsOrder, it.Name)
					}
					ret.OverrideItems = append(ret.OverrideItems, it.Name)
				}
			}
		}

		return ret, ctmiss, nil
	}

	//fmt.Printf("*** DATATYPE MISS: %s\n", b.DataType)
	if is_checkpass {
		return nil, 1, nil
	}
	return nil, 1, fmt.Errorf("Unknown data type: %s", b.DataType)
}
