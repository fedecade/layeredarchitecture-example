package httpclient

import (
	"net/http"
)

type HttpClient interface {
	PostJson(data any, extraHeaders ...HttpHeader) (*http.Response, error)
}
