package requestmanager

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
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
func (requestManager *RequestManager) Post(requesterSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, contractName string,
	functionName string, params ...interface{}) (dict.Dict, error) {
	response, err := post(false, balance.Color{}, 0, requesterSigScheme, chain, contractName, functionName, params...)
	return response, err
}

// PostWithTransfer creates a request as requester or, if not specified, as the chain originator. The contract function in the chain is called with optional params.
// It attaches 'amount' of 'color' to call. Returns response as a Dict or an error.
func (requestManager *RequestManager) PostWithTransfer(requesterSigScheme signaturescheme.SignatureScheme,
	color balance.Color, amount int64,
	chain *solo.Chain, contractName string,
	functionName string, params ...interface{}) (dict.Dict, error) {
	response, err := post(true, color, amount, requesterSigScheme, chain, contractName, functionName, params...)
	return response, err
}

func post(withTransfer bool, color balance.Color, amount int64,
	requesterSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, contractName string,
	functionName string, params ...interface{}) (dict.Dict, error) {
	request := solo.NewCallParams(contractName, functionName, params...)
	if withTransfer {
		request = request.WithTransfer(color, amount)
	}
	response, err := chain.PostRequestSync(request, requesterSigScheme)
	return response, err
}

// MustPost creates a request to contract function in the chain as requester. Fails test if request fails.
func (requestManager *RequestManager) MustPost(requesterSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, contractName string,
	functionName string, params ...interface{}) dict.Dict {
	response, err := requestManager.Post(requesterSigScheme, chain, contractName, functionName, params...)
	require.NoError(requestManager.env.T, err)
	return response
}

// MustPostWithTransfer creates a request to contract function in the chain as requester.
// It attaches 'amount' of 'color' to call. Fails test if request fails.
func (requestManager *RequestManager) MustPostWithTransfer(requesterSigScheme signaturescheme.SignatureScheme,
	color balance.Color, amount int64,
	chain *solo.Chain, contractName string,
	functionName string, params ...interface{}) dict.Dict {
	response, err := requestManager.PostWithTransfer(requesterSigScheme, color, amount, chain, contractName, functionName, params...)
	require.NoError(requestManager.env.T, err)
	return response
}

// MustPostFail creates a request to contract function in the chain as requester. Fails test if request succeeds.
func (requestManager *RequestManager) MustPostFail(requesterSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, contractName string, functionName string) {
	_, err := requestManager.Post(requesterSigScheme, chain, contractName, functionName)
	require.Error(requestManager.env.T, err)
}

// MustPostWithTransferFail creates a request to contract function in the chain as requester.
// It attaches 'amount' of 'color' to call. Fails test if request succeeds.
func (requestManager *RequestManager) MustPostWithTransferFail(requesterSigScheme signaturescheme.SignatureScheme,
	color balance.Color, amount int64,
	chain *solo.Chain, contractName string, functionName string) {
	_, err := requestManager.PostWithTransfer(requesterSigScheme, color, amount, chain, contractName, functionName)
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
