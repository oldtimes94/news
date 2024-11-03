package server

import (
	"strings"
)

type Endpoint string

const (
	Comments         Endpoint = "http://localhost:8004"
	CommentValidator Endpoint = "http://localhost:8003"
	News             Endpoint = "http://localhost:8080"
)

func (e Endpoint) Path(url string) string {
	return strings.Join([]string{string(e), url}, "")
}
