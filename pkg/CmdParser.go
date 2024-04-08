package pkg

import "flag"

type FlagParser interface {
	ParseFlagOandN(int) (bool, int, error)
}

type Parser struct {
}

func NewFlagParser() *Parser {
	return &Parser{}
}

func (fp Parser) ParseFlagOandN(defaultN int) (bool, int, error) {
	var o bool
	var n int
	flag.BoolVar(&o, "o", false, "need output")
	flag.IntVar(&n, "n", defaultN, "count of output")

	flag.Parse()
	return o, n, nil
}
