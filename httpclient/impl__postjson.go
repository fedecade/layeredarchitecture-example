package httpclient

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"example.layeredarch/logger"
	"github.com/dghubble/sling"
)

func (i *impl) PostJson(
	data any,
	extraHeaders ...HttpHeader,
) (
	*http.Response,
	error,
) {
	clen, err := i.getContentLength(data)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	handler := sling.
		New().
		Post(i.uri).
		Set("Content-Length", strconv.Itoa(clen))

	for _, h := range i.headers {
		handler.Set(h.Name, h.Value)
	}

	for _, h := range extraHeaders {
		handler.Set(h.Name, h.Value)
	}

	request, err := handler.
		BodyJSON(data).
		Request()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return i.client.Do(request)
}

func (i *impl) getContentLength(
	obj any,
) (
	int,
	error,
) {
	req, err := sling.New().BodyJSON(obj).Request()
	if err != nil {
		logger.Error(err)
		return 0, err
	}
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error(err)
		return 0, err
	}
	clen := len(content)

	return clen, nil
}
