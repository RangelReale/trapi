package trapi

func (p *Parser) parseApiHeaderList(headers []*SourceParseItemHeader, out *ApiHeaderList) error {

	for _, srcapiheader := range headers {

		// get data type
		dt, err := p.getDataType(srcapiheader.DataType)
		if err != nil {
			return err
		}

		newih := &ApiHeader{
			Name:          srcapiheader.Name,
			DataType:      dt,
			Description:   srcapiheader.Description,
			SPIB_Filename: srcapiheader.SPIB_Filename,
		}

		if _, nifound := out.List[newih.Name]; !nifound {
			out.List[newih.Name] = make([]*ApiHeader, 0)
			out.Order = append(out.Order, newih.Name)
		}
		out.List[newih.Name] = append(out.List[newih.Name], newih)
	}

	return nil
}

func (p *Parser) parseApiExampleList(example []*SourceParseItemExample, out *[]*ApiExample) {

	for _, srcapiexample := range example {

		newie := &ApiExample{
			ContentType:   srcapiexample.ContentType,
			Description:   srcapiexample.Description,
			Text:          srcapiexample.Text,
			SPIB_Filename: srcapiexample.SPIB_Filename,
		}

		*out = append(*out, newie)
	}
}
