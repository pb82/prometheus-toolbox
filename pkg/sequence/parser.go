package sequence

import (
	"errors"
	"fmt"
	"github.com/pb82/prometheus-toolbox/api"
	"github.com/pb82/prometheus-toolbox/pkg/parser"
	"strconv"
)

type Parser struct {
	parser.Parser
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

func (p *Parser) parseNullSequence() (*api.Sequence, error) {
	sequence := &api.Sequence{
		Null: true,
	}
	_, err := p.Expect(TokenTypeUnderscore)
	if err != nil {
		return nil, err
	}

	err = p.parseTimes(sequence)
	if err != nil {
		return nil, err
	}

	return sequence, nil
}

func (p *Parser) parseValueSequence() (*api.Sequence, error) {
	sequence := &api.Sequence{}

	err := p.parseInitial(sequence)
	if err != nil {
		return nil, err
	}

	err = p.parseIncrement(sequence)
	if err != nil {
		return nil, err
	}

	err = p.parseTimes(sequence)
	if err != nil {
		return nil, err
	}

	return sequence, nil
}

func (p *Parser) ParseStream() (api.SequenceGenerator, error) {
	stream := &api.Stream{}
	initial, err := p.Expect(TokenTypeNumber)
	if err != nil {
		return nil, err
	}
	initialValue, err := strconv.ParseFloat(initial.Value, 64)
	if err != nil {
		return nil, err
	}
	stream.Initial = initialValue

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
		return nil, errors.New(fmt.Sprintf(parser.ErrorUnexpectedToken, "+ or -", nextToken.Value))
	}

	increment, err := p.Expect(TokenTypeNumber)
	if err != nil {
		return nil, err
	}

	incrementValue, err := strconv.ParseFloat(increment.Value, 64)
	if err != nil {
		return nil, err
	}
	stream.Increment = incrementValue * factor
	return stream, nil
}

func (p *Parser) ParseSequence() (*api.SequenceList, error) {
	sequences := &api.SequenceList{}
	for !p.End() {
		if p.Peek().Type == TokenTypeUnderscore {
			sequence, err := p.parseNullSequence()
			if err != nil {
				return nil, err
			}

			sequences.Append(sequence)
		} else {
			sequence, err := p.parseValueSequence()
			if err != nil {
				return nil, err
			}

			sequences.Append(sequence)
		}
	}
	return sequences, nil
}

func NewParser(tokens []parser.Token) *Parser {
	return &Parser{
		Parser: parser.Parser{
			Tokens: tokens,
		},
	}
}
