package collect

type Request struct {
	ParseFunc func([]byte) ParseResult
	URL       string
	Cookie    string
	UA        string
}
type ParseResult struct {
	Requersrts []*Request
	Items      []any
}
