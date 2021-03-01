package tests

import (
	"testing"

	notsolo "github.com/brunoamancio/NotSolo"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/stretchr/testify/require"
)

func Test_MustGetInt64Result(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	const expectedDecoded = int64(1000)
	dataBytes := notSolo.Data.MustGetBytes(expectedDecoded)

	// Act
	actualDecoded := notSolo.Data.MustGetInt64(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetBytesResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	expectedDecoded := []byte{0, 0, 1}

	// Act
	actualDecoded := notSolo.Data.MustGetBytes(expectedDecoded)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetStringResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	const expectedDecoded = "test"
	dataBytes := notSolo.Data.MustGetBytes(expectedDecoded)

	// Act
	actualDecoded := notSolo.Data.MustGetString(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetHnameResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	expectedDecoded := coretypes.Hn(accounts.Interface.Name)
	dataBytes := notSolo.Data.MustGetBytes(expectedDecoded)

	// Act
	actualDecoded := notSolo.Data.MustGetHname(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetAgentIDResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	keyPair := notSolo.SigScheme.NewSignatureScheme()
	expectedAgentID := notSolo.SigScheme.MustGetAgentID(keyPair)
	dataBytes := notSolo.Data.MustGetBytes(expectedAgentID)

	// Act
	actualAgentID := notSolo.Data.MustGetAgentID(dataBytes)

	// Assert
	require.Equal(t, expectedAgentID, actualAgentID)
}

func Test_MustGetAddressResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	keyPair := notSolo.SigScheme.NewSignatureScheme()
	expectedDecoded := notSolo.SigScheme.MustGetAddress(keyPair)
	dataBytes := notSolo.Data.MustGetBytes(expectedDecoded)

	// Act
	actualDecoded := notSolo.Data.MustGetAddress(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetChainIDResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	expectedDecoded := notSolo.Chain.NewChain(nil, "dummyChain").ChainID
	dataBytes := notSolo.Data.MustGetBytes(expectedDecoded)

	// Act
	actualDecoded := notSolo.Data.MustGetChainID(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetColorResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	expectedDecoded := notSolo.Chain.NewChain(nil, "dummyChain").ChainColor
	dataBytes := notSolo.Data.MustGetBytes(expectedDecoded)

	// Act
	actualDecoded := notSolo.Data.MustGetColor(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetContractIDResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	chainName := "dummyChain"
	chain := notSolo.Chain.NewChain(nil, chainName)
	expectedDecoded := notSolo.Chain.MustGetContractID(chain, accounts.Interface.Name)
	dataBytes := notSolo.Data.MustGetBytes(expectedDecoded)

	// Act
	actualDecoded := notSolo.Data.MustGetContractID(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetHashResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	expectedDecoded := notSolo.Chain.NewChain(nil, "dummyChain").State.Hash()
	dataBytes := notSolo.Data.MustGetBytes(&expectedDecoded)

	// Act
	actualDecoded := notSolo.Data.MustGetHash(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}
