package sequence

import (
	"errors"
	"fmt"
	"promtoolbox/api"
	"promtoolbox/pkg/parser"
	"strconv"
)

type Parser struct {
	parser.Parser

	Sequences api.SequenceList
}

func (p *Parser) parseInitial(s *api.Sequence) error {
	initial, err := p.Expect(TokenTypeNumber)
	if err != nil {
		return err
	}

	initialValue, err := strconv.ParseFloat(initial.Value, 64)
	if err != nil {
		return err
	}
	s.Initial = initialValue
	return nil
}

func (p *Parser) parseIncrement(s *api.Sequence) error {
	var factor float64
	nextToken := p.Peek()
	switch nextToken.Type {
	case TokenTypePlus:
		factor = 1
		p.Consume()
	case TokenTypeMinus:
		factor = -1
		p.Consume()
	default:
		return errors.New(fmt.Sprintf(parser.ErrorUnexpectedToken, "+ or -", nextToken.Value))
	}

	increment, err := p.Expect(TokenTypeNumber)
	if err != nil {
		return err
	}

	incrementValue, err := strconv.ParseFloat(increment.Value, 64)
	if err != nil {
		return err
	}
	s.Increment = incrementValue * factor
	return nil
}

func (p *Parser) parseTimes(s *api.Sequence) error {
	_, err := p.Expect(TokenTypeX)
	if err != nil {
		return err
	}

	times, err := p.Expect(TokenTypeNumber)
	if err != nil {
		return err
	}

	timesValue, err := strconv.ParseInt(times.Value, 10, 64)
	if err != nil {
		return err
	}
	s.Times = timesValue
	return nil
}

func (p *Parser) ParseSequence() error {
	for !p.End() {
		sequence := api.Sequence{}

		err := p.parseInitial(&sequence)
		if err != nil {
			return err
		}

		err = p.parseIncrement(&sequence)
		if err != nil {
			return err
		}

		err = p.parseTimes(&sequence)
		if err != nil {
			return err
		}

		p.Sequences.Append(sequence)
	}
	return nil
}

func NewParser(tokens []parser.Token) *Parser {
	return &Parser{
		Parser: parser.Parser{
			Tokens: tokens,
		},
		Sequences: api.SequenceList{},
	}
}
