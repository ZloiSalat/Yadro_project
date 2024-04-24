package types

type Data struct {
	Comics map[int]Comics
}

type Comics struct {
	URL         string `json:"img"`
	Keywords    []string
	PreKeywords string `json:"alt"`
}
