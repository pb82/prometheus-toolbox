package sequence

import (
	"github.com/pb82/prometheus-toolbox/pkg/parser"
)

const (
	TokenTypePlus       parser.TokenType = "+"
	TokenTypeMinus      parser.TokenType = "-"
	TokenTypeX          parser.TokenType = "x"
	TokenTypeLParen     parser.TokenType = "("
	TokenTypeRParen     parser.TokenType = ")"
	TokenTypeUnderscore parser.TokenType = "_"
	TokenTypeNumber     parser.TokenType = "<number>"
	TokenTypeName       parser.TokenType = "<name>"
)
