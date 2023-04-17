package defs

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"example.layeredarch/httpclient"
	"github.com/sarulabs/di"
)

func CreditScreeningHttpClient() di.Def {
	return di.Def{
		Name:  "httpclient/creditscreening",
		Scope: di.App,
		Build: func(cnt di.Container) (interface{}, error) {
			uri := os.Getenv("CREDITSCREENING_URI")
			if len(strings.TrimSpace(uri)) == 0 {
				return nil, errors.New("CREDITSCREENING_URI is empty")
			}
			auth := os.Getenv("CREDITSCREENING_AUTH")
			if len(strings.TrimSpace(auth)) == 0 {
				return nil, errors.New("CREDITSCREENING_AUTH is empty")
			}
			return httpclient.New(
				uri,
				http.DefaultClient,
				httpclient.HttpHeader{
					Name:  "Authorization",
					Value: auth,
				},
			), nil
		},
	}
}
