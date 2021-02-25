package requestmanager

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/require"
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

// Post creates a request to the specified function in the contract in the chain as requester or, if not specified, as the chain originator. Returns response as a Dict or an error.
func (requestManager *RequestManager) Post(requesterSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, contractName string, functionName string) (dict.Dict, error) {
	request := solo.NewCallParams(contractName, functionName)
	response, err := chain.PostRequestSync(request, requesterSigScheme)
	return response, err
}

// PostMustSucceed creates a request to the specified function in the contract in the chain as requester. Fails test if request fails.
func (requestManager *RequestManager) PostMustSucceed(requesterSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, contractName string, functionName string) dict.Dict {
	response, err := requestManager.Post(requesterSigScheme, chain, contractName, functionName)
	require.NoError(requestManager.env.T, err)
	return response
}

// PostMustFail creates a request to the specified function in the contract in the chain as requester. Fails test if request succeeds.
func (requestManager *RequestManager) PostMustFail(requesterSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, contractName string, functionName string) {
	_, err := requestManager.Post(requesterSigScheme, chain, contractName, functionName)
	require.Error(requestManager.env.T, err)
}
