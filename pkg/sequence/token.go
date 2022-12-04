package sequence

import (
	"promtoolbox/pkg/parser"
)

const (
	TokenTypePlus   parser.TokenType = "+"
	TokenTypeMinus  parser.TokenType = "-"
	TokenTypeX      parser.TokenType = "x"
	TokenTypeLParen parser.TokenType = "("
	TokenTypeRParen parser.TokenType = ")"
	TokenTypeNumber parser.TokenType = "<number>"
	TokenTypeName   parser.TokenType = "<name>"
)
