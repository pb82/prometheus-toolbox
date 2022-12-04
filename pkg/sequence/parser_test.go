package sequence

import (
	"reflect"
	"testing"
	"time"
)

func TestParser(t *testing.T) {
	type testcase struct {
		input        string
		wantErr      bool
		wantSequence []int
	}

	testcases := []testcase{
		{
			input:        "0+1x10",
			wantErr:      false,
			wantSequence: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			input:        "1+1x10",
			wantErr:      false,
			wantSequence: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			input:        "10-1x10",
			wantErr:      false,
			wantSequence: []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			input:        "0+2x10",
			wantErr:      false,
			wantSequence: []int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18},
		},
		{
			input:   "0+2",
			wantErr: true,
		},
		{
			input:   "0+2x",
			wantErr: true,
		},
		{
			input:   "+2x10",
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		scanner := NewScanner(tc.input)
		scanner.Scan()

		parser := NewParser(scanner.Tokens)
		err := parser.ParseSequence()
		if err != nil && !tc.wantErr {
			t.Fatalf("\nerror: %v \n case: %v", err.Error(), tc)
		}

		seq := parser.Sequences.AsIntArray(time.Second * 1)
		if !reflect.DeepEqual(seq, tc.wantSequence) {
			t.Fatalf("\nwant: %v \n got: %v", tc.wantSequence, seq)
		}
	}
}
