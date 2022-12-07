package timeseries

import (
	"prometheus-toolbox/pkg/parser"
)

const (
	TokenTypeName   parser.TokenType = "<name>"
	TokenTypeLBrace parser.TokenType = "{"
	TokenTypeRBrace parser.TokenType = "}"
	TokenTypeEquals parser.TokenType = "="
	TokenTypeComma  parser.TokenType = ","
	TokenTypeQuote  parser.TokenType = "\""
)
