package sequence

import "github.com/pb82/prometheus-toolbox/api"

func ScanAndParseStream(source string) (*api.Stream, error) {
	scanner := NewScanner(source)
	scanner.Scan()

	parser := NewParser(scanner.Tokens)
	err := parser.ParseStream()
	if err != nil {
		return nil, err
	}
	return &parser.Stream, nil
}
