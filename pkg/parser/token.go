package parser

const (
	ErrorUnexpectedToken       = "expected '%v', but got '%v'"
	ErrorUnexpectedEndOfStream = "unexpected end of stream, expected '%v'"
)

type TokenType string

type Token struct {
	Type  TokenType
	Value string
}
