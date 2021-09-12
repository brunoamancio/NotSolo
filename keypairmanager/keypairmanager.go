package keypairmanager

import (
	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/goshimmer/packages/ledgerstate/utxodb"
	"github.com/iotaledger/hive.go/crypto/ed25519"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/colored"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/require"
)

// KeyPairManager manipulates signature structures
type KeyPairManager struct {
	env *solo.Solo
}

// New instantiates a signature scheme manager
func New(env *solo.Solo) *KeyPairManager {
	keyPairHandler := &KeyPairManager{env: env}
	return keyPairHandler
}

// NewKeyPair generates a private/public key pair. Fails test on error.
func (keyPairHandler *KeyPairManager) NewKeyPair(seed ...*ed25519.Seed) *ed25519.KeyPair {
	keyPair, address := keyPairHandler.env.NewKeyPair(seed...)
	require.NotNil(keyPairHandler.env.T, keyPair)
	require.NotNil(keyPairHandler.env.T, keyPair.PrivateKey)
	require.NotNil(keyPairHandler.env.T, keyPair.PublicKey)
	require.NotNil(keyPairHandler.env.T, address)
	keyPairHandler.RequireL1Balance(keyPair, colored.IOTA, 0)
	return keyPair
}

// NewKeyPairWithFunds generates a private/public key pair and assigns 1337 iota tokens to it (amount of funds is defined in utxodb.RequestFundsAmount)
func (keyPairHandler *KeyPairManager) NewKeyPairWithFunds(seed ...*ed25519.Seed) *ed25519.KeyPair {
	keyPair, address := keyPairHandler.env.NewKeyPairWithFunds(seed...)
	require.NotNil(keyPairHandler.env.T, keyPair)
	require.NotNil(keyPairHandler.env.T, keyPair.PrivateKey)
	require.NotNil(keyPairHandler.env.T, keyPair.PublicKey)
	require.NotNil(keyPairHandler.env.T, address)
	keyPairHandler.RequireL1Balance(keyPair, colored.IOTA, utxodb.RequestFundsAmount)
	return keyPair
}

// MustGetAgentID gets the AgentID corresponding to specified signatureScheme. Fails test on error.
func (keyPairHandler *KeyPairManager) MustGetAgentID(keyPair *ed25519.KeyPair) iscp.AgentID {
	address := ledgerstate.NewED25519Address(keyPair.PublicKey)
	agentID := *iscp.NewAgentID(address, 0)
	require.NotNil(keyPairHandler.env.T, agentID)
	return agentID
}

// MustGetAddress gets the Address (from L1) corresponding to specified signatureScheme. Fails test on error.
func (keyPairHandler *KeyPairManager) MustGetAddress(keyPair *ed25519.KeyPair) ledgerstate.Address {
	address := ledgerstate.NewED25519Address(keyPair.PublicKey)
	require.NotNil(keyPairHandler.env.T, address)
	return address
}

// RequireL1Balance verifies if the signature scheme has the expected balance of the specified color in L1.
// Fails test if balance is not equal to expectedBalance.
func (keyPairHandler *KeyPairManager) RequireL1Balance(keyPair *ed25519.KeyPair, color colored.Color, expectedBalance uint64) {
	address := ledgerstate.NewED25519Address(keyPair.PublicKey)
	keyPairHandler.env.AssertAddressBalance(address, color, expectedBalance)
}
