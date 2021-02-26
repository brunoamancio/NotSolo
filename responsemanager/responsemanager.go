package responsemanager

import (
	"github.com/iotaledger/wasp/packages/solo"
)

// ResponseManager manipulates result structures
type ResponseManager struct {
	env *solo.Solo
}

// New instantiates a result manager
func New(env *solo.Solo) *ResponseManager {
	resultHandler := &ResponseManager{env: env}
	return resultHandler
}
