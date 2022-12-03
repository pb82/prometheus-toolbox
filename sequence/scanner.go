package sequence

import (
	"promtoolbox/internal"
	"strings"
	"unicode"
)

type ScannerMode int

const (
	ScannerModeName   ScannerMode = 0
	ScannerModeNumber ScannerMode = 1
)

type Scanner struct {
	Tokens []internal.Token

	index         int
	runes         []rune
	mode          ScannerMode
	currentName   strings.Builder
	currentNumber strings.Builder
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

func (s *Scanner) commitNumber() {
	if s.currentNumber.Len() > 0 {
		s.Tokens = append(s.Tokens, internal.Token{
			Type:  TokenTypeNumber,
			Value: s.currentNumber.String(),
		})
		s.currentNumber.Reset()
	}
}

func (s *Scanner) commit() {
	switch s.mode {
	case ScannerModeName:
		s.commitName()
	case ScannerModeNumber:
		s.commitNumber()
	}
}

func (s *Scanner) append(t internal.TokenType, v ...string) {
	s.commit()

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
			s.commit()
			continue
		}

		switch nextRune {
		case 'x':
			s.append(TokenTypeX)
		case '+':
			s.append(TokenTypePlus)
		case '-':
			s.append(TokenTypeMinus)
		case '(':
			s.append(TokenTypeLParen)
		case ')':
			s.append(TokenTypeRParen)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
			s.currentNumber.WriteRune(nextRune)
			s.mode = ScannerModeNumber
		default:
			s.currentName.WriteRune(nextRune)
			s.mode = ScannerModeName
		}
	}

	s.commit()
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		Tokens:        []internal.Token{},
		index:         0,
		runes:         []rune(source),
		mode:          ScannerModeNumber,
		currentName:   strings.Builder{},
		currentNumber: strings.Builder{},
	}
}
