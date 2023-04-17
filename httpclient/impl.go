package httpclient

import "net/http"

type impl struct {
	uri     string
	client  *http.Client
	headers []HttpHeader
}

func New(
	uri string,
	client *http.Client,
	headers ...HttpHeader,
) HttpClient {
	return &impl{uri, client, headers}
}

type HttpHeader struct {
	Name  string
	Value string
}
