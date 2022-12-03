package timeseries

import (
	"promtoolbox/internal"
	"strings"
	"unicode"
)

type Scanner struct {
	Tokens []internal.Token

	index       int
	runes       []rune
	currentName strings.Builder
}

func (s *Scanner) peek() rune {
	return s.runes[s.index]
}

func (s *Scanner) next() rune {
	n := s.runes[s.index]
	s.index += 1
	return n
}

func (s *Scanner) commitName() {
	if s.currentName.Len() > 0 {
		s.Tokens = append(s.Tokens, internal.Token{
			Type:  TokenTypeName,
			Value: s.currentName.String(),
		})
		s.currentName.Reset()
	}
}

func (s *Scanner) append(t internal.TokenType, v ...string) {
	s.commitName()

	var value = string(t)
	if len(v) > 0 {
		value = strings.Join(v, "")
	}
	s.Tokens = append(s.Tokens, internal.Token{
		Type:  t,
		Value: value,
	})
}

func (s *Scanner) Scan() {
	for s.index < len(s.runes) {
		nextRune := s.next()

		if unicode.IsSpace(nextRune) {
			s.commitName()
			continue
		}

		switch nextRune {
		case '{':
			s.append(TokenTypeLBrace)
		case '}':
			s.append(TokenTypeRBrace)
		case '=':
			s.append(TokenTypeEquals)
		case ',':
			s.append(TokenTypeComma)
		case '"':
			s.append(TokenTypeQuote)
		default:
			s.currentName.WriteRune(nextRune)
		}
	}

	s.commitName()
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		Tokens:      []internal.Token{},
		index:       0,
		runes:       []rune(source),
		currentName: strings.Builder{},
	}
}
