package timeseries

import (
	"errors"
	"fmt"
	"go.buf.build/protocolbuffers/go/prometheus/prometheus"
	"prometheus-toolbox/pkg/parser"
)

type Parser struct {
	parser.Parser

	Series prometheus.TimeSeries
}

func (p *Parser) parseLabelList() error {
	_, err := p.Expect(TokenTypeLBrace)
	if err != nil {
		return err
	}

	for !p.End() {
		label := &prometheus.Label{}

		labelName, err := p.Expect(TokenTypeName)
		if err != nil {
			return err
		}
		label.Name = labelName.Value

		_, err = p.Expect(TokenTypeEquals)
		if err != nil {
			return err
		}

		_, err = p.Expect(TokenTypeQuote)
		if err != nil {
			return err
		}

		labelValue, err := p.Expect(TokenTypeName)
		if err != nil {
			return err
		}
		label.Value = labelValue.Value

		_, err = p.Expect(TokenTypeQuote)
		if err != nil {
			return err
		}

		p.Series.Labels = append(p.Series.Labels, label)

		next := p.Peek()
		if next == nil {
			break
		}
		if next.Type == TokenTypeComma {
			p.Consume()
			continue
		} else {
			break
		}
	}

	_, err = p.Expect(TokenTypeRBrace)
	return err
}

func (p *Parser) parseMetricName() error {
	nextToken, err := p.Expect(TokenTypeName)
	if err != nil {
		return err
	}

	p.Series.Labels = append(p.Series.Labels, &prometheus.Label{
		Name:  "__name__",
		Value: nextToken.Value,
	})

	if p.End() {
		return nil
	} else {
		return p.parseLabelList()
	}
}

func (p *Parser) Parse() error {
	for !p.End() {
		switch p.Peek().Type {
		case TokenTypeLBrace:
			err := p.parseLabelList()
			if err != nil {
				return err
			}
		case TokenTypeName:
			err := p.parseMetricName()
			if err != nil {
				return err
			}
		default:
			return errors.New(fmt.Sprintf(parser.ErrorUnexpectedToken,
				fmt.Sprintf("%v or %v", TokenTypeLBrace, TokenTypeName),
				p.Peek().Type))
		}
	}

	return nil
}

func NewParser(tokens []parser.Token) *Parser {
	return &Parser{
		Parser: parser.Parser{
			Tokens: tokens,
		},
		Series: prometheus.TimeSeries{},
	}
}
