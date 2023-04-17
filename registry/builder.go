package registry

import (
	"example.layeredarch/logger"
	"github.com/sarulabs/di"
)

func build(defs ...di.Def) (di.Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if err := builder.Add(defs...); err != nil {
		logger.Error(err)
		return nil, err
	}

	return builder.Build(), nil
}
