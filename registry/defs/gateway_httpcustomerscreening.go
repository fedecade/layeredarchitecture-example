package defs

import (
	"example.layeredarch/gateway/httpcreditscreenig"
	"example.layeredarch/gateway/httpcreditscreenig/translator"
	"example.layeredarch/httpclient"
	"github.com/sarulabs/di"
)

func HttpCreditScreening() di.Def {
	return di.Def{
		Name:  "gateway/httpcreditscreenig",
		Scope: di.Request,
		Build: func(cnt di.Container) (any, error) {
			client, err := cnt.SafeGet(CreditScreeningHttpClient().Name)
			if err != nil {
				return nil, err
			}
			if err != nil {
				return nil, err
			}
			trans, err := cnt.SafeGet(HttpCreditScreeningTranslator().Name)
			if err != nil {
				return nil, err
			}
			return httpcreditscreenig.New(
				client.(httpclient.HttpClient),
				trans.(translator.Translator),
			), nil
		},
	}
}
