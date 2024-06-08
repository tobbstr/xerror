package http

import (
	_ "embed"
)

//go:generate go run -tags=tools github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen -config model-config.yaml spec.yaml
