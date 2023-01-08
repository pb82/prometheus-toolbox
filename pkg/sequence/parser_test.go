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
			input:        "_x10",
			wantErr:      false,
			wantSequence: nil,
		},
		{
			input:        "1+0x2 _x10 2+0x2",
			wantErr:      false,
			wantSequence: []int{1, 1, 2, 2},
		},
		{
			input:        "_x10 1+0x2 _x10 2+0x2 _x10 3+0x2",
			wantErr:      false,
			wantSequence: []int{1, 1, 2, 2, 3, 3},
		},
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
		seq, err := parser.ParseSequence()
		if err != nil && !tc.wantErr {
			t.Fatalf("\nerror: %v \n case: %v", err.Error(), tc)
		} else if err != nil {
			continue
		}

		if !reflect.DeepEqual(seq.AsIntArray(time.Second*1), tc.wantSequence) {
			t.Fatalf("\nwant: %v \n got: %v", tc.wantSequence, seq)
		}
	}
}

func TestStreamParser(t *testing.T) {
	type testcase struct {
		input    string
		wantErr  bool
		firstTen []int64
	}

	testcases := []testcase{
		{
			input:    "0+1",
			wantErr:  false,
			firstTen: []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			input:   "0+",
			wantErr: true,
		},
		{
			input:   "+1",
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		scanner := NewScanner(tc.input)
		scanner.Scan()

		parser := NewParser(scanner.Tokens)
		stream, err := parser.ParseStream()
		if err != nil && !tc.wantErr {
			t.Fatalf("\nerror: %v \n case: %v", err.Error(), tc)
		}

		if !tc.wantErr {
			var streamValues []int64
			i := 0
			for i < 10 {
				streamValues = append(streamValues, int64(stream.Next()))
				i += 1
			}

			if !reflect.DeepEqual(streamValues, tc.firstTen) {
				t.Fatalf("\nwant: %v \n got: %v", tc.firstTen, streamValues)
			}
		}
	}
}
