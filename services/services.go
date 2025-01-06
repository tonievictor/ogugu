package service

import (
	"time"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("Service")

const dbtimeout = time.Second * 3
