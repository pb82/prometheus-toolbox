package sequence

import (
	"prometheus-toolbox/api"
)

func ScanAndParseSequence(source string) (*api.SequenceList, error) {
	scanner := NewScanner(source)
	scanner.Scan()

	parser := NewParser(scanner.Tokens)
	err := parser.ParseSequence()
	if err != nil {
		return nil, err
	}
	return &parser.Sequences, nil
}
