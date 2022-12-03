package sequence

import (
	"promtoolbox/internal"
	"reflect"
	"testing"
)

func TestScanner(t *testing.T) {
	type testcase struct {
		input      string
		wantTokens []internal.Token
	}

	testcases := []testcase{
		{
			input: "1+20x100",
			wantTokens: []internal.Token{
				{
					Type:  TokenTypeNumber,
					Value: "1",
				},
				{
					Type:  TokenTypePlus,
					Value: "+",
				},
				{
					Type:  TokenTypeNumber,
					Value: "20",
				},
				{
					Type:  TokenTypeX,
					Value: "x",
				},
				{
					Type:  TokenTypeNumber,
					Value: "100",
				},
			},
		},
		{
			input: " 1 + 2 ",
			wantTokens: []internal.Token{
				{
					Type:  TokenTypeNumber,
					Value: "1",
				},
				{
					Type:  TokenTypePlus,
					Value: "+",
				},
				{
					Type:  TokenTypeNumber,
					Value: "2",
				},
			},
		},
		{
			input: "(rnd)",
			wantTokens: []internal.Token{
				{
					Type:  TokenTypeLParen,
					Value: "(",
				},
				{
					Type:  TokenTypeName,
					Value: "rnd",
				},
				{
					Type:  TokenTypeRParen,
					Value: ")",
				},
			},
		},
		{
			input: "1 + - x 2 (rnd) 4 5",
			wantTokens: []internal.Token{
				{
					Type:  TokenTypeNumber,
					Value: "1",
				},
				{
					Type:  TokenTypePlus,
					Value: "+",
				},
				{
					Type:  TokenTypeMinus,
					Value: "-",
				},
				{
					Type:  TokenTypeX,
					Value: "x",
				},
				{
					Type:  TokenTypeNumber,
					Value: "2",
				},
				{
					Type:  TokenTypeLParen,
					Value: "(",
				},
				{
					Type:  TokenTypeName,
					Value: "rnd",
				},
				{
					Type:  TokenTypeRParen,
					Value: ")",
				},
				{
					Type:  TokenTypeNumber,
					Value: "4",
				},
				{
					Type:  TokenTypeNumber,
					Value: "5",
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
