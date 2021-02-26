package signatureschememanager

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/solo"
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

// MustGetAgentID gets the AgentID corresponding to specified signatureScheme
func (sigSchemeHandler *SignatureSchemeManager) MustGetAgentID(sigScheme signaturescheme.SignatureScheme) coretypes.AgentID {
	agentID := coretypes.NewAgentIDFromSigScheme(sigScheme)
	return agentID
}

// MustGetAddress gets the Address (from the Value Tangle) corresponding to specified signatureScheme
func (sigSchemeHandler *SignatureSchemeManager) MustGetAddress(sigScheme signaturescheme.SignatureScheme) address.Address {
	adress := sigScheme.Address()
	return adress
}

// NewSignatureScheme generates a private/public key pair
func (sigSchemeHandler *SignatureSchemeManager) NewSignatureScheme() signaturescheme.SignatureScheme {
	return sigSchemeHandler.env.NewSignatureScheme()
}

// NewSignatureSchemeWithFunds generates a private/public key pair and assigns 1337 iota tokens to it
func (sigSchemeHandler *SignatureSchemeManager) NewSignatureSchemeWithFunds() signaturescheme.SignatureScheme {
	return sigSchemeHandler.env.NewSignatureSchemeWithFunds()
}

// RequireValueTangleBalance verifies if the signature scheme has the expected balance of the specified color in the value tangle.
// Fails test if balance is not equal to expectedBalance.
func (sigSchemeHandler *SignatureSchemeManager) RequireValueTangleBalance(sigScheme signaturescheme.SignatureScheme, color balance.Color, expectedBalance int64) {
	address := sigScheme.Address()
	sigSchemeHandler.env.AssertAddressBalance(address, color, expectedBalance)
}

// RequireChainBalance verifies if the signature scheme has the expected balance of the specified color in the specified chain.
// Fails test if balance is not equal to expectedBalance.
func (sigSchemeHandler *SignatureSchemeManager) RequireChainBalance(sigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, expectedBalance int64) {
	address := sigScheme.Address()
	agentID := coretypes.NewAgentIDFromAddress(address)
	chain.AssertAccountBalance(agentID, color, expectedBalance)
}
