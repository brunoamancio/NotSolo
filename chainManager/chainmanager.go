package chainmanager

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
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
//   If 'chainOriginator' is nil, a new SignatureScheme is generated and 1337 iota tokens are assigned to it (amount of funds is defined in utxodb.RequestFundsAmount)
//   If 'validatorFeeTarget' is skipped, it is assumed equal to the chainOriginators AgentID
//
//   Fails test on error.
func (chainManager *ChainManager) NewChain(chainOriginator signaturescheme.SignatureScheme, chainName string, validatorFeeTarget ...coretypes.AgentID) *solo.Chain {
	newChain := chainManager.env.NewChain(chainOriginator, chainName, validatorFeeTarget...)
	require.NotNil(chainManager.env.T, newChain, "Could not instantiate a new chain")
	require.NotEqual(chainManager.env.T, newChain.ChainColor, balance.ColorIOTA)

	// IMPORTANT: When a chain is created, 1 IOTA is colored with the chain's color and sent to the chain's address in the value tangle
	const iotaTokensConsumedByChain = 1
	chainManager.env.AssertAddressBalance(newChain.ChainAddress, newChain.ChainColor, iotaTokensConsumedByChain)

	// IMPORTANT: When a chain is created, 1 IOTA is sent from the chain originator's account in the value tangle their account in the chain
	const iotaTokensConsumedByRequest = 1
	chainManager.RequireChainBalance(newChain.OriginatorSigScheme, newChain, balance.ColorIOTA, iotaTokensConsumedByRequest)

	// Expect zero initial fees
	feeColor, ownerFee, validatorFee := newChain.GetFeeInfo(accounts.Name)
	require.Equal(chainManager.env.T, balance.ColorIOTA, feeColor)
	require.Equal(chainManager.env.T, int64(0), ownerFee)
	require.Equal(chainManager.env.T, int64(0), validatorFee)

	chainManager.chains[chainName] = newChain
	return newChain
}

// ChangeChainFees changes the specified chains owner fee as an authorized signature scheme. See 'GrantDeployPermission' on how to (de)authorize chain changes.
// Fails test on error.
func (chainManager *ChainManager) ChangeChainFees(authorizedSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, newChainOwnerFee int64) {

	oldFeeColor, _, oldValidatorFee := chain.GetFeeInfo(accounts.Name)

	transferRequest := solo.NewCallParams(root.Interface.Name, root.FuncSetContractFee, root.ParamHname, accounts.Interface.Hname(), root.ParamOwnerFee, newChainOwnerFee)
	_, err := chain.PostRequestSync(transferRequest, chain.OriginatorSigScheme)
	require.NoError(chainManager.env.T, err)

	// Expect new fee chain owner fee
	feeColor, ownerFee, validatorFee := chain.GetFeeInfo(accounts.Name)
	require.Equal(chainManager.env.T, oldFeeColor, feeColor)
	require.Equal(chainManager.env.T, oldValidatorFee, validatorFee)
	require.Equal(chainManager.env.T, newChainOwnerFee, ownerFee)
}

// RequireChainBalance verifies if the signature scheme has the expected balance of the specified color in the specified chain.
// Fails test if balance is not equal to expectedBalance.
func (chainManager *ChainManager) RequireChainBalance(sigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, expectedBalance int64) {
	agentID := coretypes.NewAgentIDFromAddress(sigScheme.Address())
	chain.AssertAccountBalance(agentID, color, expectedBalance)
}

// GrantDeployPermission gives permission, as the chain originator, to the specified agentID to deploy SCs into the specified chain. Fails test on error.
func (chainManager *ChainManager) GrantDeployPermission(chain *solo.Chain, authorizedAgentID coretypes.AgentID) {
	err := chain.GrantDeployPermission(nil, authorizedAgentID)
	require.NoError(chainManager.env.T, err, "Could not grant deploy permission")
}

// RevokeDeployPermission revokes permission, as the chain originator, from the specified agentIDs to deploy SCs into the specified chain. Fails test on error.
func (chainManager *ChainManager) RevokeDeployPermission(chain *solo.Chain, authorizedAgentID coretypes.AgentID) {
	err := chain.RevokeDeployPermission(nil, authorizedAgentID)
	require.NoError(chainManager.env.T, err, "Could not revoke deploy permission")
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

// Transfer makes transfer of 'amount' of 'color' from the depositors account in the value tangle to the receivers account in 'chain'.
// Transfers to 'depositor' if no reciever is defined.
func (chainManager *ChainManager) Transfer(depositorSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, amount int64,
	recieverSigScheme signaturescheme.SignatureScheme) error {

	depositorAgentID := coretypes.NewAgentIDFromAddress(depositorSigScheme.Address())
	recieverAgentID := coretypes.NewAgentIDFromAddress(recieverSigScheme.Address())

	// load old balance to check against the new balance
	var params *solo.CallParams
	var oldBalance int64
	if recieverSigScheme == nil {
		params = solo.NewCallParams(accounts.Name, accounts.FuncDeposit)
		oldBalance = chain.GetAccountBalance(depositorAgentID).Balance(color)
	} else {
		params = solo.NewCallParams(accounts.Name, accounts.FuncDeposit, accounts.ParamAgentID, codec.EncodeAgentID(recieverAgentID))
		oldBalance = chain.GetAccountBalance(recieverAgentID).Balance(color)
	}

	// Transfer
	depositRequest := params.WithTransfer(color, amount)
	_, err := chain.PostRequestSync(depositRequest, depositorSigScheme)

	// load new balance to check against the old balance
	var newBalance int64
	if recieverSigScheme == nil {
		newBalance = chain.GetAccountBalance(depositorAgentID).Balance(color)
	} else {
		newBalance = chain.GetAccountBalance(recieverAgentID).Balance(color)
	}

	// checks if transfer is correct
	require.Equal(chainManager.env.T, oldBalance+amount, newBalance, "Invalid balance after transfer")
	return err
}

// TransferToSelf makes transfer of the specified amount of tokens of the specified color to the depositors account in the specified chain
func (chainManager *ChainManager) TransferToSelf(depositorSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, amount int64) error {
	err := chainManager.Transfer(depositorSigScheme, chain, color, amount, nil)
	return err
}

// MustTransferToSelf makes transfer of the specified amount of tokens of the specified color to the depositors account in the specified chain. Fails test on error.
func (chainManager *ChainManager) MustTransferToSelf(depositorSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, amount int64) {
	err := chainManager.TransferToSelf(depositorSigScheme, chain, color, amount)
	require.NoError(chainManager.env.T, err, "Could not complete transfer to self")
}

// TODO Transfer within chain (different agents)

// TODO Transfer between chains
