package pkg

import "flag"

type FlagParser interface {
	ParseFlag(int) (bool, int, error)
}

type Parser struct {
}

func NewFlagParser() *Parser {
	return &Parser{}
}

func (fp Parser) ParseFlag() (bool, error) {
	var c bool

	flag.BoolVar(&c, "c", false, "need output")

	flag.Parse()
	return c, nil
}
