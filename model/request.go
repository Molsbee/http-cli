package model

type RequestFile struct {
	Requests []Request `json:"requests"`
}

type Request struct {
	Name    string            `json:"name"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers []string          `json:"headers"`
	Data    string            `json:"data"`
	Parse   map[string]string `json:"parse"`
}
