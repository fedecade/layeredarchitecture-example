package httpcreditscreenig

import (
	"example.layeredarch/domain/creditscreening"
	"example.layeredarch/gateway/httpcreditscreenig/translator"
	"example.layeredarch/httpclient"
)

type impl struct {
	client     httpclient.HttpClient
	translator translator.Translator
}

func New(
	client httpclient.HttpClient,
	translator translator.Translator,
) creditscreening.CreditScreeing {
	return &impl{client: client, translator: translator}
}
