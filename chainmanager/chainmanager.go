package chainmanager

import (
	"errors"

	"github.com/brunoamancio/NotSolo/constants"
	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/goshimmer/packages/ledgerstate/utxodb"
	"github.com/iotaledger/hive.go/crypto/ed25519"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/colored"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
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
func (chainManager *ChainManager) NewChain(chainOriginatorKeyPair *ed25519.KeyPair, chainName string, validatorFeeTarget ...*iscp.AgentID) *solo.Chain {

	initialOriginatorBalanceInL1 := uint64(0)
	if chainOriginatorKeyPair != nil {
		chainOriginatorAddress := ledgerstate.NewED25519Address(chainOriginatorKeyPair.PublicKey)
		initialOriginatorBalanceInL1 = chainManager.env.GetAddressBalance(chainOriginatorAddress, colored.IOTA)
	}

	newChain := chainManager.env.NewChain(chainOriginatorKeyPair, chainName, validatorFeeTarget...)
	require.NotNil(chainManager.env.T, newChain, "Could not instantiate a new chain")

	// IMPORTANT: When a chain is created >>> USING SOLO <<<, a default amount of IOTA is sent to ChainID in L1
	// Another IOTA is consumed by the request and also sent to ChainID
	expectedChainIdBalance := constants.DefaultChainStartingBalance + constants.IotaTokensConsumedByRequest
	chainManager.env.AssertAddressBalance(newChain.ChainID.AsAddress(), colored.IOTA, expectedChainIdBalance)

	// IMPORTANT: Originator has no balance in the chain
	chainManager.RequireBalance(newChain.OriginatorKeyPair, newChain, colored.IOTA, 0)

	// IMPORTANT: Originator has initial balance - the amount transfered from L1
	chainOriginatorAddress := ledgerstate.NewED25519Address(newChain.OriginatorKeyPair.PublicKey)
	expectedChainOriginatorBalanceInL1 := uint64(0)
	if chainOriginatorKeyPair == nil {
		expectedChainOriginatorBalanceInL1 = utxodb.RequestFundsAmount - expectedChainIdBalance
	} else {
		expectedChainOriginatorBalanceInL1 = initialOriginatorBalanceInL1 - expectedChainIdBalance
	}

	chainManager.env.AssertAddressBalance(chainOriginatorAddress, colored.IOTA, expectedChainOriginatorBalanceInL1)

	// Expect zero initial fees
	feeColor, ownerFee, validatorFee := newChain.GetFeeInfo(accounts.Contract.Name)
	require.Equal(chainManager.env.T, colored.IOTA, feeColor)
	require.Equal(chainManager.env.T, uint64(0), ownerFee)
	require.Equal(chainManager.env.T, uint64(0), validatorFee)

	chainManager.chains[chainName] = newChain
	return newChain
}

// ChangeContractFees changes chains owner fee as 'authorized signature' scheme. Anyone with an authorized signature can use this.
// See 'GrantDeployPermission' on how to (de)authorize chain changes. Fails test on error.
func (chainManager *ChainManager) ChangeContractFees(authorizedKeyPair *ed25519.KeyPair, chain *solo.Chain, contractName string,
	newContractOwnerFee uint64) {

	oldFeeColor, _, oldValidatorFee := changeFee(chainManager, authorizedKeyPair, chain, contractName, governance.ParamOwnerFee, newContractOwnerFee)

	// Expect new fee chain owner fee
	feeColor, ownerFee, validatorFee := chain.GetFeeInfo(accounts.Contract.Name)
	require.Equal(chainManager.env.T, oldFeeColor, feeColor)
	require.Equal(chainManager.env.T, oldValidatorFee, validatorFee)
	require.Equal(chainManager.env.T, newContractOwnerFee, ownerFee)
}

// ChangeValidatorFees changes the validator fee as 'authorized signature' scheme. Anyone with an authorized signature can use this.
// See 'GrantDeployPermission' on how to (de)authorize chain changes. Fails test on error.
func (chainManager *ChainManager) ChangeValidatorFees(authorizedKeyPair *ed25519.KeyPair, chain *solo.Chain, contractName string,
	newValidatorFee uint64) {
	oldFeeColor, oldOwnerFee, _ := changeFee(chainManager, authorizedKeyPair, chain, contractName, governance.ParamValidatorFee, newValidatorFee)

	// Expect new fee chain owner fee
	feeColor, ownerFee, validatorFee := chain.GetFeeInfo(accounts.Contract.Name)
	require.Equal(chainManager.env.T, oldFeeColor, feeColor)
	require.Equal(chainManager.env.T, newValidatorFee, validatorFee)
	require.Equal(chainManager.env.T, oldOwnerFee, ownerFee)
}

func changeFee(chainManager *ChainManager, authorizedKeyPair *ed25519.KeyPair, chain *solo.Chain, contractName string,
	feeParam string, newFee uint64) (oldFeeColor colored.Color, oldChainOwnerFee uint64, oldValidatorFee uint64) {

	contractRecord, err := chainManager.GetContractRecord(chain, contractName)
	require.NoError(chainManager.env.T, err)
	require.NotNil(chainManager.env.T, contractRecord, "Contract could not be found")

	oldFeeColor, oldChainOwnerFee, oldValidatorFee = chain.GetFeeInfo(contractName)

	request := solo.NewCallParams(governance.Contract.Name, governance.FuncSetContractFee.Name, governance.ParamHname, contractRecord.Hname(), feeParam, newFee).WithIotas(constants.IotaTokensConsumedByRequest)
	_, err = chain.PostRequestSync(request, authorizedKeyPair)
	require.NoError(chainManager.env.T, err)

	return oldFeeColor, oldChainOwnerFee, oldValidatorFee
}

// GrantDeployPermission gives permission, as the chain originator, to 'authorizedKeyPair' to deploy SCs into the specified chain. Fails test on error.
func (chainManager *ChainManager) GrantDeployPermission(chain *solo.Chain, authorizedKeyPair *ed25519.KeyPair) {
	authorizedAddress := ledgerstate.NewED25519Address(authorizedKeyPair.PublicKey)
	authorizedAgentID := iscp.NewAgentID(authorizedAddress, 0)
	err := chain.GrantDeployPermission(nil, *authorizedAgentID)
	require.NoError(chainManager.env.T, err, "Could not grant deploy permission")
}

// RevokeDeployPermission revokes permission, as the chain originator, from 'authorizedKeyPair' to deploy SCs into 'chain'. Fails test on error.
func (chainManager *ChainManager) RevokeDeployPermission(chain *solo.Chain, authorizedKeyPair *ed25519.KeyPair) {
	authorizedAddress := ledgerstate.NewED25519Address(authorizedKeyPair.PublicKey)
	authorizedAgentID := iscp.NewAgentID(authorizedAddress, 0)
	err := chain.RevokeDeployPermission(nil, *authorizedAgentID)
	require.NoError(chainManager.env.T, err, "Could not revoke deploy permission")
}

// GrantAgentDeployPermission gives permission, as the chain originator, to 'authorizedAgentID' to deploy SCs into the specified chain. Fails test on error.
func (chainManager *ChainManager) GrantAgentDeployPermission(chain *solo.Chain, authorizedAgentID iscp.AgentID) {
	err := chain.GrantDeployPermission(nil, authorizedAgentID)
	require.NoError(chainManager.env.T, err, "Could not grant deploy permission")
}

// RevokeAgentDeployPermission revokes permission, as the chain originator, from 'authorizedAgentID' to deploy SCs into 'chain'. Fails test on error.
func (chainManager *ChainManager) RevokeAgentDeployPermission(chain *solo.Chain, authorizedAgentID iscp.AgentID) {
	err := chain.RevokeDeployPermission(nil, authorizedAgentID)
	require.NoError(chainManager.env.T, err, "Could not revoke deploy permission")
}

// MustGetAgentID ensures 'chain' contains 'contract' and returns its ContractID. Fails test on error.
func (chainManager *ChainManager) MustGetAgentID(chain *solo.Chain, contractName string) *iscp.AgentID {
	chain, ok := chainManager.chains[chain.Name]
	require.True(chainManager.env.T, ok)
	require.NotNil(chainManager.env.T, chain)

	contractID := chain.ContractAgentID(contractName)
	return contractID
}

// DeployWasmContract uploads and deploys 'constract wasm file'
func (chainManager *ChainManager) DeployWasmContract(chain *solo.Chain, contractOriginatorKeyPair *ed25519.KeyPair, contractName string, contractWasmFilePath string) {
	err := chain.DeployWasmContract(contractOriginatorKeyPair, contractName, contractWasmFilePath)
	require.NoError(chainManager.env.T, err, "Could not deploy wasm contract")
}

// NewChainAndDeployWasmContract calls NewChain and then DeployWasmContract
func (chainManager *ChainManager) NewChainAndDeployWasmContract(chainOriginatorKeyPair *ed25519.KeyPair, chainName string,
	contractOriginatorKeyPair *ed25519.KeyPair, contractName string, contractWasmFilePath string,
	validatorFeeTarget ...*iscp.AgentID) (*solo.Chain, *root.ContractRecord) {

	chain := chainManager.NewChain(chainOriginatorKeyPair, chainName, validatorFeeTarget...)
	chainManager.DeployWasmContract(chain, contractOriginatorKeyPair, contractName, contractWasmFilePath)
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

// MustTransferToL1ToSelf makes transfer of 'amount' of 'color' from the depositors account in 'chain' to the depositors address in L1.
// Fails test on error.
// Important: Due to the current logic in the accounts contract (from IOTA Foundation), all tokens of 'depositor' are withdrawn from 'chain' to his address in L1.
func (chainManager *ChainManager) MustTransferToL1ToSelf(depositorKeyPair *ed25519.KeyPair, chain *solo.Chain, color colored.Color, transferAmount uint64) {
	err := chainManager.TransferToL1ToSelf(depositorKeyPair, chain, color, transferAmount)
	require.NoError(chainManager.env.T, err)
}

// TransferToL1ToSelf makes transfer of 'amount' of 'color' from the depositors account in 'chain' to the depositors address in L1.
// Important: Due to the current logic in the accounts contract (from IOTA Foundation), all tokens of 'depositor' are withdrawn from 'chain' to his address in L1.
func (chainManager *ChainManager) TransferToL1ToSelf(depositorKeyPair *ed25519.KeyPair, chain *solo.Chain, color colored.Color, transferAmount uint64) error {

	request := solo.NewCallParams(accounts.Contract.Name, accounts.FuncWithdraw.Name).WithIotas(transferAmount)
	_, err := chain.PostRequestSync(request, depositorKeyPair)

	return err
}

// MustTransferBetweenChains makes transfer of 'amount' of 'color' from the depositors account in 'sourceChain' to the receivers account in 'destinationChain'.
// Transfers to 'depositor' if no receiver is defined. Fails test on error.
// Important: Due to the current logic in the accounts contract (from IOTA Foundation), all tokens of 'depositor' are withdrawn from 'sourceChain' to his address in L1.
func (chainManager *ChainManager) MustTransferBetweenChains(depositorKeyPair *ed25519.KeyPair, sourceChain *solo.Chain, color colored.Color, transferAmount uint64,
	destinationChain *solo.Chain, receiverKeyPair *ed25519.KeyPair) {
	err := chainManager.TransferBetweenChains(depositorKeyPair, sourceChain, color, transferAmount, destinationChain, receiverKeyPair)
	require.NoError(chainManager.env.T, err)
}

// TransferBetweenChains makes transfer of 'amount' of 'color' from the depositors account in 'sourceChain' to the receivers account in 'destinationChain'.
// Transfers to 'depositor' if no receiver is defined.
// Important: Due to the current logic in the accounts contract (from IOTA Foundation), all tokens of 'depositor' are withdrawn from 'sourceChain' to his address in L1.
func (chainManager *ChainManager) TransferBetweenChains(depositorKeyPair *ed25519.KeyPair, sourceChain *solo.Chain, color colored.Color, transferAmount uint64,
	destinationChain *solo.Chain, receiverKeyPair *ed25519.KeyPair) error {

	isReceiverDefined := receiverKeyPair != nil

	if !isReceiverDefined {
		receiverKeyPair = depositorKeyPair
	}

	// Transfer from 'sourceChain' to depositor's account in L1
	err := chainManager.TransferToL1ToSelf(depositorKeyPair, sourceChain, color, transferAmount)

	if err != nil {
		return err
	}

	// Transfer from depositor's address in L1 to receiver's account in 'destinationChain'
	receiverAddress := ledgerstate.NewED25519Address(receiverKeyPair.PublicKey)
	receiverAgentID := iscp.NewAgentID(receiverAddress, 0)

	request := solo.NewCallParams(accounts.Contract.Name, accounts.FuncDeposit.Name, accounts.ParamAgentID, receiverAgentID).
		WithTransfer(color, transferAmount)
	_, err = destinationChain.PostRequestSync(request, depositorKeyPair)

	return err
}

// TransferWithinChain makes transfer of 'amount' of 'color' from the depositors account in 'chain' to the receivers account in the same chain.
// Nothing is transfered if no receiver is defined.
// Important: Due to the current logic in the accounts contract (from IOTA Foundation), all tokens are withdrawn from 'sourceChain', not only the specified color and amount.
func (chainManager *ChainManager) TransferWithinChain(depositorKeyPair *ed25519.KeyPair, chain *solo.Chain, color colored.Color, transferAmount uint64,
	receiverKeyPair *ed25519.KeyPair) error {
	isReceiverDefined := receiverKeyPair != nil

	if !isReceiverDefined {
		return errors.New("receiver not defined")
	}

	// Transfer from depositor's account in 'chain' to the receiver's account in the same chain.
	err := chainManager.TransferBetweenChains(depositorKeyPair, chain, color, transferAmount, chain, receiverKeyPair)
	return err
}

// MustTransferWithinChain makes transfer of 'amount' of 'color' from the depositors account in 'chain' to the receivers account in the same chain.
// Nothing is transfered if no receiver is defined. Fails test on error.
// Important: Due to the current logic in the accounts contract (from IOTA Foundation), all tokens are withdrawn from 'sourceChain', not only the specified color and amount.
func (chainManager *ChainManager) MustTransferWithinChain(depositorKeyPair *ed25519.KeyPair, chain *solo.Chain, color colored.Color, transferAmount uint64,
	receiverKeyPair *ed25519.KeyPair) {
	err := chainManager.TransferWithinChain(depositorKeyPair, chain, color, transferAmount, receiverKeyPair)
	require.NoError(chainManager.env.T, err)
}

// RequireBalance verifies if the signature scheme has the expected balance of 'color' in 'chain'.
// Fails test if balance is not equal to expectedBalance.
func (chainManager *ChainManager) RequireBalance(keyPair *ed25519.KeyPair, chain *solo.Chain, color colored.Color, expectedBalance uint64) {
	address := ledgerstate.NewED25519Address(keyPair.PublicKey)
	agentID := iscp.NewAgentID(address, 0)

	chain.AssertAccountBalance(agentID, color, expectedBalance)
}

// RequireContractBalance verifies if 'contract' has the expected balance of 'color' in 'chain'.
// Fails test if contract is neither defined, nor found, or if the balance is not equal to expectedBalance.
func (chainManager *ChainManager) RequireContractBalance(chain *solo.Chain, contractName string, color colored.Color, expectedBalance uint64) {

	// Get contract record
	contractRecord, err := chain.FindContract(contractName)
	require.NoError(chainManager.env.T, err)
	require.NotNil(chainManager.env.T, contractRecord, "Contract could not be found")

	// Get contract's AgentID
	contractAgentID := chain.ContractAgentID(contractRecord.Name)

	chain.AssertAccountBalance(contractAgentID, color, expectedBalance)
}
