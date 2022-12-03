package internal

import (
	"errors"
	"fmt"
)

type Parser struct {
	index  int
	Tokens []Token
}

func (p *Parser) End() bool {
	return p.index >= len(p.Tokens)
}

func (p *Parser) Consume() {
	p.index += 1
}

func (p *Parser) Peek() *Token {
	if p.End() {
		return nil
	}
	return &p.Tokens[p.index]
}

func (p *Parser) Expect(tokenType TokenType) (*Token, error) {
	if p.End() {
		return nil, errors.New(fmt.Sprintf(ErrorUnexpectedEndOfStream, tokenType))
	}
	nextType := p.Tokens[p.index].Type
	if nextType == tokenType {
		nextToken := &p.Tokens[p.index]
		p.index += 1
		return nextToken, nil
	}
	return nil, errors.New(fmt.Sprintf(ErrorUnexpectedToken, tokenType, nextType))
}
