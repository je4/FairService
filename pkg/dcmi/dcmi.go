package dcmi

type DCMI struct {
	OAI_DCNS          string `xml:"xmlns:oai_dc,attr"`
	DCNS              string `xml:"xmlns:dc"`
	XsiType           string `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string `xml:"xsi:schemaLocation,attr"`
	Title             string `xml:"dc:title"`

	Type DCMIType `xml:"dc:type"`
}

func (dcmi *DCMI) InitNamespace() {
	dcmi.OAI_DCNS = "http://www.openarchives.org/OAI/2.0/oai_dc/"
	dcmi.DCNS = "http://purl.org/dc/elements/1.1/"
	dcmi.XsiType = "http://www.w3.org/2001/XMLSchema-instance"
	dcmi.XsiSchemaLocation = "http://www.openarchives.org/OAI/2.0/oai_dc/ http://www.openarchives.org/OAI/2.0/oai_dc.xsd"
}
