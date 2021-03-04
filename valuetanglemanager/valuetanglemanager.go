package valuetanglemanager

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/stretchr/testify/require"
)

// ValueTangleManager manipulates chains.
type ValueTangleManager struct {
	env *solo.Solo
}

// New instantiates a chain manager.
func New(env *solo.Solo) *ValueTangleManager {
	valueTangleManager := &ValueTangleManager{env: env}
	return valueTangleManager
}

// MustTransferToChain makes transfer of 'amount' of 'color' from the depositors account in the value tangle to the receivers account in 'chain'.
// Transfers to 'depositor' if no receiver is defined.
// Fails test on error.
func (valueTangleManager *ValueTangleManager) MustTransferToChain(depositorSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, transferAmount int64,
	receiverSigScheme signaturescheme.SignatureScheme) {
	valueTangleManager.TransferToChain(depositorSigScheme, chain, color, transferAmount, receiverSigScheme)
}

// TransferToChain makes transfer of 'amount' of 'color' from the depositors account in the value tangle to the receivers account in 'chain'.
// Transfers to 'depositor' if no receiver is defined.
func (valueTangleManager *ValueTangleManager) TransferToChain(depositorSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, transferAmount int64,
	receiverSigScheme signaturescheme.SignatureScheme) error {

	isReceiverDefined := receiverSigScheme != nil

	if !isReceiverDefined {
		receiverSigScheme = depositorSigScheme
	}

	receiverAddress := receiverSigScheme.Address()
	receiverAgentID := coretypes.NewAgentIDFromAddress(receiverAddress)

	// Transfer
	err := transferToAgent(depositorSigScheme, chain, color, transferAmount, receiverAgentID)
	return err
}

// MustTransferToChainToSelf makes transfer of 'amount' of 'color' from the depositors account in the value tangle to the depositors account in 'chain'. Fails test on error.
func (valueTangleManager *ValueTangleManager) MustTransferToChainToSelf(depositorSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, amount int64) {
	err := valueTangleManager.TransferToChainToSelf(depositorSigScheme, chain, color, amount)
	require.NoError(valueTangleManager.env.T, err, "Could not complete transfer to self")
}

// TransferToChainToSelf makes transfer of 'amount' of 'color' from the depositors account in the value tangle to the depositors account in 'chain'.
func (valueTangleManager *ValueTangleManager) TransferToChainToSelf(depositorSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, amount int64) error {
	err := valueTangleManager.TransferToChain(depositorSigScheme, chain, color, amount, nil)
	return err
}

// MustTransferToContract makes transfer of 'amount' of 'color' from the depositors account in the value tangle to the contract's account in 'chain'.
// Nothing is transfered if no contract is neither defined nor found. Fails test on error.
func (valueTangleManager *ValueTangleManager) MustTransferToContract(depositorSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, transferAmount int64,
	contractName string) {

	err := valueTangleManager.TransferToContract(depositorSigScheme, chain, color, transferAmount, contractName)
	require.NoError(valueTangleManager.env.T, err)
}

// TransferToContract makes transfer of 'amount' of 'color' from the depositors account in the value tangle to the contract's account in 'chain'.
// Nothing is transfered if no contract is neither defined nor found.
func (valueTangleManager *ValueTangleManager) TransferToContract(depositorSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, transferAmount int64,
	contractName string) error {

	isContractDefined := contractName != ""
	if !isContractDefined {
		return nil
	}

	// Get contract record
	contractRecord, err := chain.FindContract(contractName)
	require.NoError(valueTangleManager.env.T, err)

	if contractRecord == nil {
		return nil
	}

	// Get contract's AgentID
	contractID := coretypes.NewContractID(chain.ChainID, contractRecord.Hname())
	contractAgentID := coretypes.NewAgentIDFromContractID(contractID)

	// Transfer
	err = transferToAgent(depositorSigScheme, chain, color, transferAmount, contractAgentID)
	return err
}

func transferToAgent(depositorSigScheme signaturescheme.SignatureScheme, chain *solo.Chain, color balance.Color, transferAmount int64,
	agentID coretypes.AgentID) error {

	params := solo.NewCallParams(accounts.Name, accounts.FuncDeposit, accounts.ParamAgentID, codec.EncodeAgentID(agentID))
	depositRequest := params.WithTransfer(color, transferAmount)
	_, err := chain.PostRequestSync(depositRequest, depositorSigScheme)
	return err
}

// RequireBalance verifies if the signature scheme has the expected balance of 'color' in the value tangle.
// Fails test if balance is not equal to expectedBalance.
func (valueTangleManager *ValueTangleManager) RequireBalance(sigScheme signaturescheme.SignatureScheme, color balance.Color, expectedBalance int64) {
	valueTangleManager.RequireAddressBalance(sigScheme.Address(), color, expectedBalance)
}

// RequireAddressBalance verifies if the address has the expected balance of 'color' in the value tangle.
// Fails test if balance is not equal to expectedBalance.
func (valueTangleManager *ValueTangleManager) RequireAddressBalance(address address.Address, color balance.Color, expectedBalance int64) {
	valueTangleManager.env.AssertAddressBalance(address, color, expectedBalance)
}
