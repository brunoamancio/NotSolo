package signatureschememanager

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/goshimmer/dapps/waspconn/packages/utxodb"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/require"
)

// SignatureSchemeManager manipulates signature structures
type SignatureSchemeManager struct {
	env *solo.Solo
}

// New instantiates a signature scheme manager
func New(env *solo.Solo) *SignatureSchemeManager {
	signatureSchemeHandler := &SignatureSchemeManager{env: env}
	return signatureSchemeHandler
}

// NewSignatureScheme generates a private/public key pair. Fails test on error.
func (sigSchemeHandler *SignatureSchemeManager) NewSignatureScheme() signaturescheme.SignatureScheme {
	sigScheme := sigSchemeHandler.env.NewSignatureScheme()
	require.NotNil(sigSchemeHandler.env.T, sigScheme)
	sigSchemeHandler.RequireValueTangleBalance(sigScheme, balance.ColorIOTA, 0)
	return sigScheme
}

// NewSignatureSchemeWithFunds generates a private/public key pair and assigns 1337 iota tokens to it (amount of funds is defined in utxodb.RequestFundsAmount)
func (sigSchemeHandler *SignatureSchemeManager) NewSignatureSchemeWithFunds() signaturescheme.SignatureScheme {
	sigScheme := sigSchemeHandler.env.NewSignatureSchemeWithFunds()
	require.NotNil(sigSchemeHandler.env.T, sigScheme)
	sigSchemeHandler.RequireValueTangleBalance(sigScheme, balance.ColorIOTA, utxodb.RequestFundsAmount)
	return sigScheme
}

// MustGetAgentID gets the AgentID corresponding to specified signatureScheme. Fails test on error.
func (sigSchemeHandler *SignatureSchemeManager) MustGetAgentID(sigScheme signaturescheme.SignatureScheme) coretypes.AgentID {
	agentID := coretypes.NewAgentIDFromSigScheme(sigScheme)
	require.NotNil(sigSchemeHandler.env.T, agentID)
	return agentID
}

// MustGetAddress gets the Address (from the Value Tangle) corresponding to specified signatureScheme. Fails test on error.
func (sigSchemeHandler *SignatureSchemeManager) MustGetAddress(sigScheme signaturescheme.SignatureScheme) address.Address {
	address := sigScheme.Address()
	require.NotNil(sigSchemeHandler.env.T, address)
	return address
}

// RequireValueTangleBalance verifies if the signature scheme has the expected balance of the specified color in the value tangle.
// Fails test if balance is not equal to expectedBalance.
func (sigSchemeHandler *SignatureSchemeManager) RequireValueTangleBalance(sigScheme signaturescheme.SignatureScheme, color balance.Color, expectedBalance int64) {
	address := sigScheme.Address()
	sigSchemeHandler.env.AssertAddressBalance(address, color, expectedBalance)
}
