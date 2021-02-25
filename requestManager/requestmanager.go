package requestmanager

import "testing"

// RequestManager manipulates requests
type RequestManager struct {
	t *testing.T
}

// New instantiates a request manager
func New(t *testing.T) *RequestManager {
	requestManager := &RequestManager{t: t}
	return requestManager
}
