package requestmanager

import (
	"github.com/iotaledger/wasp/packages/solo"
)

// RequestManager manipulates requests
type RequestManager struct {
	env *solo.Solo
}

// New instantiates a request manager
func New(env *solo.Solo) *RequestManager {
	requestManager := &RequestManager{env: env}
	return requestManager
}
