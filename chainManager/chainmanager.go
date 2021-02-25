package chainmanager

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/require"
)

// ChainManager manipulates chains
type ChainManager struct {
	env    *solo.Solo
	chains map[string]*solo.Chain
}

// New instantiates a chain manager
func New(env *solo.Solo) *ChainManager {
	chainManager := &ChainManager{env: env, chains: make(map[string]*solo.Chain)}
	return chainManager
}

// NewChain instantiates a new chain
//   If 'chainOriginator' is nil, a new SignatureScheme is generated and 1337 iota tokens are assigned to it
//   If 'validatorFeeTarget' is skipped, it is assumed equal to the chainOriginator's AgentID
func (chainManager *ChainManager) NewChain(chainOriginator signaturescheme.SignatureScheme, chainName string, validatorFeeTarget ...coretypes.AgentID) *solo.Chain {
	newChain := chainManager.env.NewChain(chainOriginator, chainName, validatorFeeTarget...)
	chainManager.chains[chainName] = newChain
	return newChain
}

// MustGetContractID ensures the specified chain contains the specified contract and returns its ContractID
func (chainManager *ChainManager) MustGetContractID(chainName string, contractName string) coretypes.ContractID {
	chain, ok := chainManager.chains[chainName]
	require.True(chainManager.env.T, ok)
	require.NotNil(chainManager.env.T, chain)

	contractID := coretypes.NewContractID(chain.ChainID, coretypes.Hn(contractName))
	return contractID
}
