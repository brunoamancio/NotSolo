package tests

import (
	"testing"

	notsolo "github.com/brunoamancio/NotSolo"
	"github.com/brunoamancio/NotSolo/constants"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/goshimmer/dapps/waspconn/packages/utxodb"
)

func Test_TransferToChainToSelf(t *testing.T) {
	notSolo := notsolo.New(t)

	// Create a chain
	chain := notSolo.Chain.NewChain(nil, "myChain")

	// Create a sigscheme with dummy funds (amount is defined in utxodb.RequestFundsAmount)
	senderSigScheme := notSolo.SigScheme.NewSignatureSchemeWithFunds()
	transferAmount := int64(100)

	// Send some funds to chain
	notSolo.ValueTangle.MustTransferToChainToSelf(senderSigScheme, chain, balance.ColorIOTA, transferAmount)
	notSolo.ValueTangle.RequireBalance(senderSigScheme, balance.ColorIOTA, utxodb.RequestFundsAmount-transferAmount-constants.IotaTokensConsumedByRequest)
	notSolo.Chain.RequireBalance(senderSigScheme, chain, balance.ColorIOTA, transferAmount+constants.IotaTokensConsumedByRequest)
}

func Test_TransferToValueTangleToSelf(t *testing.T) {
	notSolo := notsolo.New(t)

	// Create a chain
	chain := notSolo.Chain.NewChain(nil, "myChain")

	// Create a sigscheme with dummy funds (amount is defined in utxodb.RequestFundsAmount)
	senderSigScheme := notSolo.SigScheme.NewSignatureSchemeWithFunds()
	transferAmount := int64(100)

	// Send some funds to chain
	notSolo.ValueTangle.MustTransferToChainToSelf(senderSigScheme, chain, balance.ColorIOTA, transferAmount)
	balanceInValueTangle := utxodb.RequestFundsAmount - transferAmount - constants.IotaTokensConsumedByRequest
	balanceInChain := transferAmount + constants.IotaTokensConsumedByRequest
	notSolo.ValueTangle.RequireBalance(senderSigScheme, balance.ColorIOTA, balanceInValueTangle)
	notSolo.Chain.RequireBalance(senderSigScheme, chain, balance.ColorIOTA, balanceInChain)

	// Send funds from chain to value tangle
	notSolo.Chain.MustTransferToValueTangleToSelf(senderSigScheme, chain, balance.ColorIOTA, balanceInChain)
	notSolo.Chain.RequireBalance(senderSigScheme, chain, balance.ColorIOTA, 0)
	notSolo.ValueTangle.RequireBalance(senderSigScheme, balance.ColorIOTA, balanceInValueTangle+balanceInChain)
}

func Test_TransferBetweenChains(t *testing.T) {
	notSolo := notsolo.New(t)

	// Create sourceChain and destinationChain
	sourceChain := notSolo.Chain.NewChain(nil, "mySourceChain")
	destinationChain := notSolo.Chain.NewChain(nil, "myDestinationChain")

	// Create a sigscheme with dummy funds (amount is defined in utxodb.RequestFundsAmount)
	senderSigScheme := notSolo.SigScheme.NewSignatureSchemeWithFunds()
	transferAmount := int64(100)
	receiverSigScheme := notSolo.SigScheme.NewSignatureScheme()

	// Send some funds to chain
	notSolo.ValueTangle.MustTransferToChainToSelf(senderSigScheme, sourceChain, balance.ColorIOTA, transferAmount)
	senderBalanceInValueTangle := utxodb.RequestFundsAmount - transferAmount - constants.IotaTokensConsumedByRequest
	senderBalanceInChain := transferAmount + constants.IotaTokensConsumedByRequest
	notSolo.ValueTangle.RequireBalance(senderSigScheme, balance.ColorIOTA, senderBalanceInValueTangle)
	notSolo.Chain.RequireBalance(senderSigScheme, sourceChain, balance.ColorIOTA, senderBalanceInChain)

	// Send funds from sourceChain to destinationChain
	notSolo.Chain.MustTransferBetweenChains(senderSigScheme, sourceChain, balance.ColorIOTA, senderBalanceInChain, destinationChain, receiverSigScheme)
	notSolo.ValueTangle.RequireBalance(senderSigScheme, balance.ColorIOTA, senderBalanceInValueTangle-constants.IotaTokensConsumedByRequest)
	notSolo.Chain.RequireBalance(senderSigScheme, sourceChain, balance.ColorIOTA, 0)
	notSolo.Chain.RequireBalance(senderSigScheme, destinationChain, balance.ColorIOTA, constants.IotaTokensConsumedByRequest)
	notSolo.Chain.RequireBalance(receiverSigScheme, destinationChain, balance.ColorIOTA, senderBalanceInChain)
}

func Test_TransferWithinChain(t *testing.T) {
	notSolo := notsolo.New(t)

	// Create sourceChain and destinationChain
	chain := notSolo.Chain.NewChain(nil, "mySourceChain")

	// Create a sigscheme with dummy funds (amount is defined in utxodb.RequestFundsAmount)
	senderSigScheme := notSolo.SigScheme.NewSignatureSchemeWithFunds()
	transferAmount := int64(100)
	receiverSigScheme := notSolo.SigScheme.NewSignatureScheme()

	// Send some funds to chain
	notSolo.ValueTangle.MustTransferToChainToSelf(senderSigScheme, chain, balance.ColorIOTA, transferAmount)
	senderBalanceInValueTangle := utxodb.RequestFundsAmount - transferAmount - constants.IotaTokensConsumedByRequest
	senderBalanceInChain := transferAmount + constants.IotaTokensConsumedByRequest
	notSolo.ValueTangle.RequireBalance(senderSigScheme, balance.ColorIOTA, senderBalanceInValueTangle)
	notSolo.Chain.RequireBalance(senderSigScheme, chain, balance.ColorIOTA, senderBalanceInChain)

	// Send funds from sourceChain to destinationChain
	notSolo.Chain.MustTransferWithinChain(senderSigScheme, chain, balance.ColorIOTA, senderBalanceInChain, receiverSigScheme)
	notSolo.ValueTangle.RequireBalance(senderSigScheme, balance.ColorIOTA, senderBalanceInValueTangle-constants.IotaTokensConsumedByRequest)
	notSolo.Chain.RequireBalance(senderSigScheme, chain, balance.ColorIOTA, constants.IotaTokensConsumedByRequest)
	notSolo.Chain.RequireBalance(receiverSigScheme, chain, balance.ColorIOTA, senderBalanceInChain)
}
