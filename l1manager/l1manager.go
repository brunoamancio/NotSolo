package l1manager

import (
	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/hive.go/crypto/ed25519"

	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/colored"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/stretchr/testify/require"
)

// L1Manager manipulates chains.
type L1Manager struct {
	env *solo.Solo
}

// New instantiates a chain manager.
func New(env *solo.Solo) *L1Manager {
	l1Manager := &L1Manager{env: env}
	return l1Manager
}

// MustTransferToChain makes transfer of 'amount' of 'color' from the depositors account in L1 to the receivers account in 'chain'.
// Transfers to 'depositor' if no receiver is defined.
// Fails test on error.
func (l1Manager *L1Manager) MustTransferToChain(depositorKeyPair *ed25519.KeyPair, chain *solo.Chain, color colored.Color, transferAmount uint64,
	receiverKeyPair *ed25519.KeyPair) {
	l1Manager.TransferToChain(depositorKeyPair, chain, color, transferAmount, receiverKeyPair)
}

// TransferToChain makes transfer of 'amount' of 'color' from the depositors account in L1 to the receivers account in 'chain'.
// Transfers to 'depositor' if no receiver is defined.
func (l1Manager *L1Manager) TransferToChain(depositorKeyPair *ed25519.KeyPair, chain *solo.Chain, color colored.Color, transferAmount uint64,
	receiverKeyPair *ed25519.KeyPair) error {

	isReceiverDefined := receiverKeyPair != nil

	if !isReceiverDefined {
		receiverKeyPair = depositorKeyPair
	}
	receiverAddress := ledgerstate.NewED25519Address(receiverKeyPair.PublicKey)
	receiverAgentID := iscp.NewAgentID(receiverAddress, 0)

	// Transfer
	err := transferToAgent(depositorKeyPair, chain, color, transferAmount, receiverAgentID)
	return err
}

// MustTransferToChainToSelf makes transfer of 'amount' of 'color' from the depositors account in L1 to the depositors account in 'chain'. Fails test on error.
func (l1Manager *L1Manager) MustTransferToChainToSelf(depositorKeyPair *ed25519.KeyPair, chain *solo.Chain, color colored.Color, amount uint64) {
	err := l1Manager.TransferToChainToSelf(depositorKeyPair, chain, color, amount)
	require.NoError(l1Manager.env.T, err, "Could not complete transfer to self")
}

// TransferToChainToSelf makes transfer of 'amount' of 'color' from the depositors account in L1 to the depositors account in 'chain'.
func (l1Manager *L1Manager) TransferToChainToSelf(depositorKeyPair *ed25519.KeyPair, chain *solo.Chain, color colored.Color, amount uint64) error {
	err := l1Manager.TransferToChain(depositorKeyPair, chain, color, amount, nil)
	return err
}

// MustTransferToContract makes transfer of 'amount' of 'color' from the depositors account in L1 to the contract's account in 'chain'.
// Nothing is transfered if no contract is neither defined nor found. Fails test on error.
func (l1Manager *L1Manager) MustTransferToContract(depositorKeyPair *ed25519.KeyPair, chain *solo.Chain, color colored.Color, transferAmount uint64,
	contractName string) {

	err := l1Manager.TransferToContract(depositorKeyPair, chain, color, transferAmount, contractName)
	require.NoError(l1Manager.env.T, err)
}

// TransferToContract makes transfer of 'amount' of 'color' from the depositors account in L1 to the contract's account in 'chain'.
// Nothing is transfered if no contract is neither defined nor found.
func (l1Manager *L1Manager) TransferToContract(depositorKeyPair *ed25519.KeyPair, chain *solo.Chain, color colored.Color, transferAmount uint64,
	contractName string) error {

	isContractDefined := contractName != ""
	if !isContractDefined {
		return nil
	}

	// Get contract record
	contractRecord, err := chain.FindContract(contractName)
	require.NoError(l1Manager.env.T, err)

	if contractRecord == nil {
		return nil
	}

	// Get contract's AgentID
	contractAgentID := iscp.NewAgentID(chain.ChainID.AsAddress(), contractRecord.Hname())

	// Transfer
	err = transferToAgent(depositorKeyPair, chain, color, transferAmount, contractAgentID)
	return err
}

func transferToAgent(depositorKeypair *ed25519.KeyPair, chain *solo.Chain, color colored.Color, transferAmount uint64,
	agentID *iscp.AgentID) error {

	params := solo.NewCallParams(accounts.Contract.Name, accounts.FuncDeposit.Name, accounts.ParamAgentID, codec.EncodeAgentID(agentID))
	depositRequest := params.WithTransfer(color, transferAmount)
	_, err := chain.PostRequestSync(depositRequest, depositorKeypair)
	return err
}

// RequireBalance verifies if the signature scheme has the expected balance of 'color' in L1.
// Fails test if balance is not equal to expectedBalance.
func (l1Manager *L1Manager) RequireBalance(keyPair *ed25519.KeyPair, color colored.Color, expectedBalance uint64) {
	address := ledgerstate.NewED25519Address(keyPair.PublicKey)
	l1Manager.RequireAddressBalance(address, color, expectedBalance)
}

// RequireAddressBalance verifies if the address has the expected balance of 'color' in L1.
// Fails test if balance is not equal to expectedBalance.
func (l1Manager *L1Manager) RequireAddressBalance(address ledgerstate.Address, color colored.Color, expectedBalance uint64) {
	l1Manager.env.AssertAddressBalance(address, color, expectedBalance)
}
