package sequence

import (
	"promtoolbox/internal"
)

const (
	TokenTypePlus   internal.TokenType = "+"
	TokenTypeMinus  internal.TokenType = "-"
	TokenTypeX      internal.TokenType = "x"
	TokenTypeLParen internal.TokenType = "("
	TokenTypeRParen internal.TokenType = ")"
	TokenTypeNumber internal.TokenType = "<number>"
	TokenTypeName   internal.TokenType = "<name>"
)
