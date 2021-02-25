package signatureschememanager

import (
	"testing"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/solo"
)

// SignatureSchemeManager manipulates signature structures
type SignatureSchemeManager struct {
	t   *testing.T
	env *solo.Solo
}

// New instantiates a signature scheme manager
func New(env *solo.Solo) *SignatureSchemeManager {
	signatureSchemeHandler := &SignatureSchemeManager{t: env.T, env: env}
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
