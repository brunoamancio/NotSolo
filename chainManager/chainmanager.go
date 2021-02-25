package chainmanager

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/stretchr/testify/require"
)

// ChainManager manipulates chains
type ChainManager struct {
	env    *solo.Solo
	chains map[string]*solo.Chain
}

// Dispose implements Disposable for ChainManager
func (chainManager *ChainManager) Dispose() {
	chainManager.chains = make(map[string]*solo.Chain)
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

// MustGetContractID ensures the specified chain contains the specified contract and returns its ContractID. Fails test on error.
func (chainManager *ChainManager) MustGetContractID(chainName string, contractName string) coretypes.ContractID {
	chain, ok := chainManager.chains[chainName]
	require.True(chainManager.env.T, ok)
	require.NotNil(chainManager.env.T, chain)

	contractID := coretypes.NewContractID(chain.ChainID, coretypes.Hn(contractName))
	return contractID
}

// DeployWasmContract uploads and deploys the specified constract wasm file
func (chainManager *ChainManager) DeployWasmContract(chain *solo.Chain, contractOriginator signaturescheme.SignatureScheme, contractName string, contractWasmFilePath string) {
	err := chain.DeployWasmContract(contractOriginator, contractName, contractWasmFilePath)
	require.NoError(chainManager.env.T, err, "Could not deploy wasm contract")
}

// NewChainAndDeployWasmContract calls NewChain and then DeployWasmContract
func (chainManager *ChainManager) NewChainAndDeployWasmContract(chainOriginator signaturescheme.SignatureScheme, chainName string,
	contractOriginator signaturescheme.SignatureScheme, contractName string, contractWasmFilePath string,
	validatorFeeTarget ...coretypes.AgentID) (*solo.Chain, *root.ContractRecord) {

	chain := chainManager.NewChain(chainOriginator, chainName, validatorFeeTarget...)
	chainManager.DeployWasmContract(chain, contractOriginator, contractName, contractWasmFilePath)
	contractRecord := chainManager.MustGetContractRecord(chain, contractName)
	return chain, contractRecord
}

// GetContractRecord searches the specified chain for the specified contract record by name
func (chainManager *ChainManager) GetContractRecord(chain *solo.Chain, contractName string) (*root.ContractRecord, error) {
	return chain.FindContract(contractName)
}

// MustGetContractRecord searches the specified chain for the specified contract record by name. Fails test on error.
func (chainManager *ChainManager) MustGetContractRecord(chain *solo.Chain, contractName string) *root.ContractRecord {
	contractRecord, err := chainManager.GetContractRecord(chain, contractName)
	require.NoError(chainManager.env.T, err, "Could not find contract")
	require.NotNil(chainManager.env.T, contractRecord)
	return contractRecord
}
