package responsemanager

import (
	"testing"
)

// ResponseManager manipulates result structures
type ResponseManager struct {
	t *testing.T
}

// New instantiates a result manager
func New(t *testing.T) *ResponseManager {
	resultHandler := &ResponseManager{t: t}
	return resultHandler
}
