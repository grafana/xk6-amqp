package amqp

import (
	"go.k6.io/k6/js/modules"
)

const version = "v0.0.1"

type Amqp struct {
    Version string
}

func init() {
	modules.Register("k6/x/amqp", &Amqp{
		Version: version,
	})
}
