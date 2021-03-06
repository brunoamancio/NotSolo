package chainmanager

import (
	"errors"

	"github.com/brunoamancio/NotSolo/constants"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/packages/coretypes"
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
	chainManager.env.AssertAddressBalance(newChain.ChainAddress, newChain.ChainColor, constants.IotaTokensConsumedByChain)

	// IMPORTANT: When a chain is created, 1 IOTA is sent from the chain originator's account in the value tangle their account in the chain
	chainManager.RequireBalance(newChain.OriginatorSigScheme, newChain, balance.ColorIOTA, constants.IotaTokensConsumedByRequest)

	// Expect zero initial fees
	feeColor, ownerFee, validatorFee := newChain.GetFeeInfo(accounts.Name)
	require.Equal(chainManager.env.T, balance.ColorIOTA, feeColor)
	require.Equal(chainManager.env.T, int64(0), ownerFee)
	require.Equal(chainManager.env.T, int64(0), validatorFee)

	chainManager.chains[chainName] = newChain
	return newChain
}

// ChangeContractFees changes chains owner fee as 'authorized signature' scheme. Anyone with an authorized signature can use this.
// See 'GrantDeployPermission' on how to (de)authorize chain changes. Fails test on error.
func (chainManager *ChainManager) ChangeContractFees(authorizedSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, contractName string,
	newContractOwnerFee int64) {

	oldFeeColor, _, oldValidatorFee := changeFee(chainManager, authorizedSigScheme, chain, contractName, root.ParamOwnerFee, newContractOwnerFee)

	// Expect new fee chain owner fee
	feeColor, ownerFee, validatorFee := chain.GetFeeInfo(accounts.Name)
	require.Equal(chainManager.env.T, oldFeeColor, feeColor)
	require.Equal(chainManager.env.T, oldValidatorFee, validatorFee)
	require.Equal(chainManager.env.T, newContractOwnerFee, ownerFee)
}

// ChangeValidatorFees changes the validator fee as 'authorized signature' scheme. Anyone with an authorized signature can use this.
// See 'GrantDeployPermission' on how to (de)authorize chain changes. Fails test on error.
func (chainManager *ChainManager) ChangeValidatorFees(authorizedSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, contractName string,
	newValidatorFee int64) {
	oldFeeColor, oldOwnerFee, _ := changeFee(chainManager, authorizedSigScheme, chain, contractName, root.ParamValidatorFee, newValidatorFee)

	// Expect new fee chain owner fee
	feeColor, ownerFee, validatorFee := chain.GetFeeInfo(accounts.Name)
	require.Equal(chainManager.env.T, oldFeeColor, feeColor)
	require.Equal(chainManager.env.T, newValidatorFee, validatorFee)
	require.Equal(chainManager.env.T, oldOwnerFee, ownerFee)
}

func changeFee(chainManager *ChainManager, authorizedSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, contractName string,
	feeParam string, newFee int64) (oldFeeColor balance.Color, oldChainOwnerFee int64, oldValidatorFee int64) {

	contractRecord, err := chainManager.GetContractRecord(chain, contractName)
	require.NoError(chainManager.env.T, err)
	require.NotNil(chainManager.env.T, contractRecord, "Contract could not be found")

	oldFeeColor, oldChainOwnerFee, oldValidatorFee = chain.GetFeeInfo(contractName)

	request := solo.NewCallParams(root.Interface.Name, root.FuncSetContractFee, root.ParamHname, contractRecord.Hname(), feeParam, newFee)
	_, err = chain.PostRequestSync(request, authorizedSigScheme)
	require.NoError(chainManager.env.T, err)

	return oldFeeColor, oldChainOwnerFee, oldValidatorFee
}

// GrantDeployPermission gives permission, as the chain originator, to 'authorizedSigScheme' to deploy SCs into the specified chain. Fails test on error.
func (chainManager *ChainManager) GrantDeployPermission(chain *solo.Chain, authorizedSigScheme signaturescheme.SignatureScheme) {
	authorizedAgentID := coretypes.NewAgentIDFromSigScheme(authorizedSigScheme)
	err := chain.GrantDeployPermission(nil, authorizedAgentID)
	require.NoError(chainManager.env.T, err, "Could not grant deploy permission")
}

// RevokeDeployPermission revokes permission, as the chain originator, from 'authorizedSigScheme' to deploy SCs into 'chain'. Fails test on error.
func (chainManager *ChainManager) RevokeDeployPermission(chain *solo.Chain, authorizedSigScheme signaturescheme.SignatureScheme) {
	authorizedAgentID := coretypes.NewAgentIDFromSigScheme(authorizedSigScheme)
	err := chain.RevokeDeployPermission(nil, authorizedAgentID)
	require.NoError(chainManager.env.T, err, "Could not revoke deploy permission")
}

// GrantAgentDeployPermission gives permission, as the chain originator, to 'authorizedAgentID' to deploy SCs into the specified chain. Fails test on error.
func (chainManager *ChainManager) GrantAgentDeployPermission(chain *solo.Chain, authorizedAgentID coretypes.AgentID) {
	err := chain.GrantDeployPermission(nil, authorizedAgentID)
	require.NoError(chainManager.env.T, err, "Could not grant deploy permission")
}

// RevokeAgentDeployPermission revokes permission, as the chain originator, from 'authorizedAgentID' to deploy SCs into 'chain'. Fails test on error.
func (chainManager *ChainManager) RevokeAgentDeployPermission(chain *solo.Chain, authorizedAgentID coretypes.AgentID) {
	err := chain.RevokeDeployPermission(nil, authorizedAgentID)
	require.NoError(chainManager.env.T, err, "Could not revoke deploy permission")
}

// MustGetContractID ensures 'chain' contains 'contract' and returns its ContractID. Fails test on error.
func (chainManager *ChainManager) MustGetContractID(chain *solo.Chain, contractName string) coretypes.ContractID {
	chain, ok := chainManager.chains[chain.Name]
	require.True(chainManager.env.T, ok)
	require.NotNil(chainManager.env.T, chain)

	contractID := coretypes.NewContractID(chain.ChainID, coretypes.Hn(contractName))
	return contractID
}

// DeployWasmContract uploads and deploys 'constract wasm file'
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

// GetContractRecord searches 'chain' for 'contract record' by name
func (chainManager *ChainManager) GetContractRecord(chain *solo.Chain, contractName string) (*root.ContractRecord, error) {
	return chain.FindContract(contractName)
}

// MustGetContractRecord searches 'chain' for 'contract record' by name. Fails test on error.
func (chainManager *ChainManager) MustGetContractRecord(chain *solo.Chain, contractName string) *root.ContractRecord {
	contractRecord, err := chainManager.GetContractRecord(chain, contractName)
	require.NoError(chainManager.env.T, err, "Could not find contract")
	require.NotNil(chainManager.env.T, contractRecord)
	return contractRecord
}

// MustTransferToValueTangleToSelf makes transfer of 'amount' of 'color' from the depositors account in 'chain' to the depositors address in the value tangle.
// Fails test on error.
// Important: Due to the current logic in the accounts contract (from IOTA Foundation), all tokens of 'depositor' are withdrawn from 'chain' to his address in the value tangle.
func (chainManager *ChainManager) MustTransferToValueTangleToSelf(depositorSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, transferAmount int64) {
	err := chainManager.TransferToValueTangleToSelf(depositorSigScheme, chain, color, transferAmount)
	require.NoError(chainManager.env.T, err)
}

// TransferToValueTangleToSelf makes transfer of 'amount' of 'color' from the depositors account in 'chain' to the depositors address in the value tangle.
// Important: Due to the current logic in the accounts contract (from IOTA Foundation), all tokens of 'depositor' are withdrawn from 'chain' to his address in the value tangle.
func (chainManager *ChainManager) TransferToValueTangleToSelf(depositorSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, transferAmount int64) error {

	request := solo.NewCallParams(accounts.Interface.Name, accounts.FuncWithdrawToAddress)
	_, err := chain.PostRequestSync(request, depositorSigScheme)

	return err
}

// MustTransferBetweenChains makes transfer of 'amount' of 'color' from the depositors account in 'sourceChain' to the receivers account in 'destinationChain'.
// Transfers to 'depositor' if no receiver is defined. Fails test on error.
// Important: Due to the current logic in the accounts contract (from IOTA Foundation), all tokens of 'depositor' are withdrawn from 'sourceChain' to his address in the value tangle.
func (chainManager *ChainManager) MustTransferBetweenChains(depositorSigScheme signaturescheme.SignatureScheme, sourceChain *solo.Chain, color balance.Color, transferAmount int64,
	destinationChain *solo.Chain, receiverSigScheme signaturescheme.SignatureScheme) {
	err := chainManager.TransferBetweenChains(depositorSigScheme, sourceChain, color, transferAmount, destinationChain, receiverSigScheme)
	require.NoError(chainManager.env.T, err)
}

// TransferBetweenChains makes transfer of 'amount' of 'color' from the depositors account in 'sourceChain' to the receivers account in 'destinationChain'.
// Transfers to 'depositor' if no receiver is defined.
// Important: Due to the current logic in the accounts contract (from IOTA Foundation), all tokens of 'depositor' are withdrawn from 'sourceChain' to his address in the value tangle.
func (chainManager *ChainManager) TransferBetweenChains(depositorSigScheme signaturescheme.SignatureScheme, sourceChain *solo.Chain, color balance.Color, transferAmount int64,
	destinationChain *solo.Chain, receiverSigScheme signaturescheme.SignatureScheme) error {

	isReceiverDefined := receiverSigScheme != nil

	if !isReceiverDefined {
		receiverSigScheme = depositorSigScheme
	}

	// Transfer from 'sourceChain' to depositor's account in the value tangle
	err := chainManager.TransferToValueTangleToSelf(depositorSigScheme, sourceChain, color, transferAmount)

	if err != nil {
		return err
	}

	// Transfer from depositor's address in the value tangle to receiver's account in 'destinationChain'
	receiverAgentID := coretypes.NewAgentIDFromSigScheme(receiverSigScheme)
	request := solo.NewCallParams(accounts.Name, accounts.FuncDeposit, accounts.ParamAgentID, receiverAgentID).
		WithTransfer(color, transferAmount)
	_, err = destinationChain.PostRequestSync(request, depositorSigScheme)

	return err
}

// TransferWithinChain makes transfer of 'amount' of 'color' from the depositors account in 'chain' to the receivers account in the same chain.
// Nothing is transfered if no receiver is defined.
// Important: Due to the current logic in the accounts contract (from IOTA Foundation), all tokens are withdrawn from 'sourceChain', not only the specified color and amount.
func (chainManager *ChainManager) TransferWithinChain(depositorSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, transferAmount int64,
	receiverSigScheme signaturescheme.SignatureScheme) error {
	isReceiverDefined := receiverSigScheme != nil

	if !isReceiverDefined {
		return errors.New("Receiver not defined")
	}

	// Transfer from depositor's account in 'chain' to the receiver's account in the same chain.
	err := chainManager.TransferBetweenChains(depositorSigScheme, chain, color, transferAmount, chain, receiverSigScheme)
	return err
}

// MustTransferWithinChain makes transfer of 'amount' of 'color' from the depositors account in 'chain' to the receivers account in the same chain.
// Nothing is transfered if no receiver is defined. Fails test on error.
// Important: Due to the current logic in the accounts contract (from IOTA Foundation), all tokens are withdrawn from 'sourceChain', not only the specified color and amount.
func (chainManager *ChainManager) MustTransferWithinChain(depositorSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, transferAmount int64,
	receiverSigScheme signaturescheme.SignatureScheme) {
	err := chainManager.TransferWithinChain(depositorSigScheme, chain, color, transferAmount, receiverSigScheme)
	require.NoError(chainManager.env.T, err)
}

// RequireBalance verifies if the signature scheme has the expected balance of 'color' in 'chain'.
// Fails test if balance is not equal to expectedBalance.
func (chainManager *ChainManager) RequireBalance(sigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, expectedBalance int64) {
	agentID := coretypes.NewAgentIDFromSigScheme(sigScheme)
	chain.AssertAccountBalance(agentID, color, expectedBalance)
}

// RequireContractBalance verifies if 'contract' has the expected balance of 'color' in 'chain'.
// Fails test if contract is neither defined, nor found, or if the balance is not equal to expectedBalance.
func (chainManager *ChainManager) RequireContractBalance(chain *solo.Chain, contractName string, color balance.Color, expectedBalance int64) {

	// Get contract record
	contractRecord, err := chain.FindContract(contractName)
	require.NoError(chainManager.env.T, err)
	require.NotNil(chainManager.env.T, contractRecord, "Contract could not be found")

	// Get contract's AgentID
	contractID := coretypes.NewContractID(chain.ChainID, contractRecord.Hname())
	contractAgentID := coretypes.NewAgentIDFromContractID(contractID)

	chain.AssertAccountBalance(contractAgentID, color, expectedBalance)
}
