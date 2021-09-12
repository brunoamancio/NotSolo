package tests

import (
	"testing"

	notsolo "github.com/brunoamancio/NotSolo"
	"github.com/iotaledger/goshimmer/packages/ledgerstate/utxodb"
	"github.com/iotaledger/wasp/packages/iscp/colored"
)

func Test_TransferToChainToSelf(t *testing.T) {
	notSolo := notsolo.New(t)

	// Create a chain
	chain := notSolo.Chain.NewChain(nil, "myChain")

	// Create a key pair with dummy funds (amount is defined in utxodb.RequestFundsAmount)
	senderKeyPair := notSolo.KeyPair.NewKeyPairWithFunds()
	transferAmount := uint64(100)

	// Send some funds to chain
	notSolo.L1.MustTransferToChainToSelf(senderKeyPair, chain, colored.IOTA, transferAmount)
	notSolo.L1.RequireBalance(senderKeyPair, colored.IOTA, utxodb.RequestFundsAmount-transferAmount)
	notSolo.Chain.RequireBalance(senderKeyPair, chain, colored.IOTA, transferAmount)
}

func Test_TransferToL1ToSelf(t *testing.T) {
	notSolo := notsolo.New(t)

	// Create a chain
	chain := notSolo.Chain.NewChain(nil, "myChain")

	// Create a key pair with dummy funds (amount is defined in utxodb.RequestFundsAmount)
	senderKeyPair := notSolo.KeyPair.NewKeyPairWithFunds()
	transferAmount := uint64(100)

	// Send some funds to chain
	notSolo.L1.MustTransferToChainToSelf(senderKeyPair, chain, colored.IOTA, transferAmount)
	balanceInL1 := utxodb.RequestFundsAmount - transferAmount
	balanceInChain := transferAmount
	notSolo.L1.RequireBalance(senderKeyPair, colored.IOTA, balanceInL1)
	notSolo.Chain.RequireBalance(senderKeyPair, chain, colored.IOTA, balanceInChain)

	// Send funds from chain to L1
	notSolo.Chain.MustTransferToL1ToSelf(senderKeyPair, chain, colored.IOTA, balanceInChain)
	notSolo.Chain.RequireBalance(senderKeyPair, chain, colored.IOTA, 0)
	notSolo.L1.RequireBalance(senderKeyPair, colored.IOTA, balanceInL1+balanceInChain)
}

func Test_TransferBetweenChains(t *testing.T) {
	notSolo := notsolo.New(t)

	// Create sourceChain and destinationChain
	sourceChain := notSolo.Chain.NewChain(nil, "mySourceChain")
	destinationChain := notSolo.Chain.NewChain(nil, "myDestinationChain")

	// Create a key pair with dummy funds (amount is defined in utxodb.RequestFundsAmount)
	senderKeyPair := notSolo.KeyPair.NewKeyPairWithFunds()
	transferAmount := uint64(100)
	receiverKeyPair := notSolo.KeyPair.NewKeyPair()

	// Send some funds to chain
	notSolo.L1.MustTransferToChainToSelf(senderKeyPair, sourceChain, colored.IOTA, transferAmount)
	senderBalanceInL1 := utxodb.RequestFundsAmount - transferAmount
	senderBalanceInChain := transferAmount
	notSolo.L1.RequireBalance(senderKeyPair, colored.IOTA, senderBalanceInL1)
	notSolo.Chain.RequireBalance(senderKeyPair, sourceChain, colored.IOTA, senderBalanceInChain)

	// Send funds from sourceChain to destinationChain
	notSolo.Chain.MustTransferBetweenChains(senderKeyPair, sourceChain, colored.IOTA, senderBalanceInChain, destinationChain, receiverKeyPair)
	notSolo.L1.RequireBalance(senderKeyPair, colored.IOTA, senderBalanceInL1)
	notSolo.Chain.RequireBalance(senderKeyPair, sourceChain, colored.IOTA, 0)
	notSolo.Chain.RequireBalance(senderKeyPair, destinationChain, colored.IOTA, 0)
	notSolo.Chain.RequireBalance(receiverKeyPair, destinationChain, colored.IOTA, senderBalanceInChain)
}

func Test_TransferWithinChain(t *testing.T) {
	notSolo := notsolo.New(t)

	// Create a chain
	chain := notSolo.Chain.NewChain(nil, "mySourceChain")

	// Create a key pair with dummy funds (amount is defined in utxodb.RequestFundsAmount)
	senderKeyPair := notSolo.KeyPair.NewKeyPairWithFunds()
	transferAmount := uint64(100)
	receiverKeyPair := notSolo.KeyPair.NewKeyPair()

	// Send some funds to chain
	notSolo.L1.MustTransferToChainToSelf(senderKeyPair, chain, colored.IOTA, transferAmount)
	senderBalanceInL1 := utxodb.RequestFundsAmount - transferAmount
	senderBalanceInChain := transferAmount
	notSolo.L1.RequireBalance(senderKeyPair, colored.IOTA, senderBalanceInL1)
	notSolo.Chain.RequireBalance(senderKeyPair, chain, colored.IOTA, senderBalanceInChain)

	// Send funds from sourceChain to destinationChain
	notSolo.Chain.MustTransferWithinChain(senderKeyPair, chain, colored.IOTA, senderBalanceInChain, receiverKeyPair)
	notSolo.L1.RequireBalance(senderKeyPair, colored.IOTA, senderBalanceInL1)
	notSolo.Chain.RequireBalance(senderKeyPair, chain, colored.IOTA, 0)
	notSolo.Chain.RequireBalance(receiverKeyPair, chain, colored.IOTA, senderBalanceInChain)
}
