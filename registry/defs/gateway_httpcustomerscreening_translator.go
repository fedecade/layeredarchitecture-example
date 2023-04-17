package defs

import (
	"example.layeredarch/gateway/httpcreditscreenig/translator"
	"github.com/sarulabs/di"
)

func HttpCreditScreeningTranslator() di.Def {
	return di.Def{
		Name:  "gateway/httpcreditscreenig/translator",
		Scope: di.App,
		Build: func(cnt di.Container) (any, error) {
			return translator.New(), nil
		},
	}
}
