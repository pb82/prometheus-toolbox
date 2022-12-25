package sequence

import "github.com/pb82/prometheus-toolbox/api"

func ScanAndParseStream(source string) (*api.Stream, error) {
	scanner := NewScanner(source)
	scanner.Scan()

	parser := NewParser(scanner.Tokens)
	stream, err := parser.ParseStream()
	if err != nil {
		return nil, err
	}
	return stream, nil
}
