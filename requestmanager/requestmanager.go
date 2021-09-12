package requestmanager

import (
	"github.com/iotaledger/hive.go/crypto/ed25519"
	"github.com/iotaledger/wasp/packages/iscp/colored"
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

// Post creates a request as requester or, if not specified, as the chain originator. The contract function in the chain is called with optional params.
// Returns response as a Dict or an error.
func (requestManager *RequestManager) Post(requesterKeyPair *ed25519.KeyPair, chain *solo.Chain, contractName string,
	functionName string, params ...interface{}) (dict.Dict, error) {
	response, err := post(false, colored.Color{}, 0, requesterKeyPair, chain, contractName, functionName, params...)
	return response, err
}

// PostWithTransfer creates a request as requester or, if not specified, as the chain originator. The contract function in the chain is called with optional params.
// It attaches 'amount' of 'color' to call. Returns response as a Dict or an error.
func (requestManager *RequestManager) PostWithTransfer(requesterKeyPair *ed25519.KeyPair,
	color colored.Color, amount uint64,
	chain *solo.Chain, contractName string,
	functionName string, params ...interface{}) (dict.Dict, error) {
	response, err := post(true, color, amount, requesterKeyPair, chain, contractName, functionName, params...)
	return response, err
}

func post(withTransfer bool, color colored.Color, amount uint64,
	requesterKeyPair *ed25519.KeyPair, chain *solo.Chain, contractName string,
	functionName string, params ...interface{}) (dict.Dict, error) {
	request := solo.NewCallParams(contractName, functionName, params...)
	if withTransfer {
		request = request.WithTransfer(color, amount)
	}
	response, err := chain.PostRequestSync(request, requesterKeyPair)
	return response, err
}

// MustPost creates a request to contract function in the chain as requester. Fails test if request fails.
func (requestManager *RequestManager) MustPost(requesterKeyPair *ed25519.KeyPair, chain *solo.Chain, contractName string,
	functionName string, params ...interface{}) dict.Dict {
	response, err := requestManager.Post(requesterKeyPair, chain, contractName, functionName, params...)
	require.NoError(requestManager.env.T, err)
	return response
}

// MustPostWithTransfer creates a request to contract function in the chain as requester.
// It attaches 'amount' of 'color' to call. Fails test if request fails.
func (requestManager *RequestManager) MustPostWithTransfer(requesterKeyPair *ed25519.KeyPair,
	color colored.Color, amount uint64,
	chain *solo.Chain, contractName string,
	functionName string, params ...interface{}) dict.Dict {
	response, err := requestManager.PostWithTransfer(requesterKeyPair, color, amount, chain, contractName, functionName, params...)
	require.NoError(requestManager.env.T, err)
	return response
}

// MustPostFail creates a request to contract function in the chain as requester. Fails test if request succeeds.
func (requestManager *RequestManager) MustPostFail(requesterKeyPair *ed25519.KeyPair, chain *solo.Chain, contractName string, functionName string) {
	_, err := requestManager.Post(requesterKeyPair, chain, contractName, functionName)
	require.Error(requestManager.env.T, err)
}

// MustPostWithTransferFail creates a request to contract function in the chain as requester.
// It attaches 'amount' of 'color' to call. Fails test if request succeeds.
func (requestManager *RequestManager) MustPostWithTransferFail(requesterKeyPair *ed25519.KeyPair,
	color colored.Color, amount uint64,
	chain *solo.Chain, contractName string, functionName string) {
	_, err := requestManager.PostWithTransfer(requesterKeyPair, color, amount, chain, contractName, functionName)
	require.Error(requestManager.env.T, err)
}

// View creates a view request. The contract view in the chain is called with optional params.
// Returns response as a Dict or an error.
func (requestManager *RequestManager) View(chain *solo.Chain, contractName string,
	functionName string, params ...interface{}) (dict.Dict, error) {
	response, err := chain.CallView(contractName, functionName, params...)
	return response, err
}

// MustView creates a view request. The contract view in the chain is called with optional params.
// Returns response as a Dict. Fails test on error.
func (requestManager *RequestManager) MustView(chain *solo.Chain, contractName string,
	functionName string, params ...interface{}) dict.Dict {
	response, err := chain.CallView(contractName, functionName, params...)
	require.NoError(requestManager.env.T, err)
	return response
}

// MustViewFail creates a view request. The contract view in the chain is called with optional params.
// Fails test if request succeeds.
func (requestManager *RequestManager) MustViewFail(chain *solo.Chain, contractName string,
	functionName string, params ...interface{}) {
	_, err := chain.CallView(contractName, functionName, params...)
	require.Error(requestManager.env.T, err)
}
