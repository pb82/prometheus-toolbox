package sequence

import (
	"github.com/pb82/prometheus-toolbox/api"
)

func ScanAndParseSequence(source string) (api.SequenceGenerator, error) {
	scanner := NewScanner(source)
	scanner.Scan()

	parser := NewParser(scanner.Tokens)
	sequences, err := parser.ParseSequence()
	if err != nil {
		return nil, err
	}
	return sequences, nil
}
