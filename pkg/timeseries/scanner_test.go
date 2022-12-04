package timeseries

import (
	"promtoolbox/pkg/parser"
	"reflect"
	"testing"
)

func TestScanner(t *testing.T) {
	type testcase struct {
		input      string
		wantTokens []parser.Token
	}

	testcases := []testcase{
		{
			input: "name",
			wantTokens: []parser.Token{
				{
					Type:  TokenTypeName,
					Value: "name",
				},
			},
		},
		{
			input: "{}",
			wantTokens: []parser.Token{
				{
					Type:  TokenTypeLBrace,
					Value: "{",
				},
				{
					Type:  TokenTypeRBrace,
					Value: "}",
				},
			},
		},
		{
			input: "metric{label=\"value\"}",
			wantTokens: []parser.Token{
				{
					Type:  TokenTypeName,
					Value: "metric",
				},
				{
					Type:  TokenTypeLBrace,
					Value: "{",
				},
				{
					Type:  TokenTypeName,
					Value: "label",
				},
				{
					Type:  TokenTypeEquals,
					Value: "=",
				},
				{
					Type:  TokenTypeQuote,
					Value: "\"",
				},
				{
					Type:  TokenTypeName,
					Value: "value",
				},
				{
					Type:  TokenTypeQuote,
					Value: "\"",
				},
				{
					Type:  TokenTypeRBrace,
					Value: "}",
				},
			},
		},
	}

	for _, tc := range testcases {
		scanner := NewScanner(tc.input)
		scanner.Scan()
		if !reflect.DeepEqual(scanner.Tokens, tc.wantTokens) {
			t.Fatalf("\nwant: %v \n got: %v", tc.wantTokens, scanner.Tokens)
		}
	}
}
