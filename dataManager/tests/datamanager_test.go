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
	dataBytes := notSolo.DataManager.MustGetBytes(expectedDecoded)

	// Act
	actualDecoded := notSolo.DataManager.MustGetInt64(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetBytesResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	expectedDecoded := []byte{0, 0, 1}

	// Act
	actualDecoded := notSolo.DataManager.MustGetBytes(expectedDecoded)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetStringResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	const expectedDecoded = "test"
	dataBytes := notSolo.DataManager.MustGetBytes(expectedDecoded)

	// Act
	actualDecoded := notSolo.DataManager.MustGetString(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetHnameResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	expectedDecoded := coretypes.Hn(accounts.Interface.Name)
	dataBytes := notSolo.DataManager.MustGetBytes(expectedDecoded)

	// Act
	actualDecoded := notSolo.DataManager.MustGetHname(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetAgentIDResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	keyPair := notSolo.SignatureSchemeManager.NewSignatureScheme()
	expectedAgentID := notSolo.SignatureSchemeManager.MustGetAgentID(keyPair)
	dataBytes := notSolo.DataManager.MustGetBytes(expectedAgentID)

	// Act
	actualAgentID := notSolo.DataManager.MustGetAgentID(dataBytes)

	// Assert
	require.Equal(t, expectedAgentID, actualAgentID)
}

func Test_MustGetAddressResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	keyPair := notSolo.SignatureSchemeManager.NewSignatureScheme()
	expectedDecoded := notSolo.SignatureSchemeManager.MustGetAddress(keyPair)
	dataBytes := notSolo.DataManager.MustGetBytes(expectedDecoded)

	// Act
	actualDecoded := notSolo.DataManager.MustGetAddress(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetChainIDResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	expectedDecoded := notSolo.ChainManager.NewChain(nil, "dummyChain").ChainID
	dataBytes := notSolo.DataManager.MustGetBytes(expectedDecoded)

	// Act
	actualDecoded := notSolo.DataManager.MustGetChainID(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetColorResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	expectedDecoded := notSolo.ChainManager.NewChain(nil, "dummyChain").ChainColor
	dataBytes := notSolo.DataManager.MustGetBytes(expectedDecoded)

	// Act
	actualDecoded := notSolo.DataManager.MustGetColor(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetContractIDResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	chainName := "dummyChain"
	notSolo.ChainManager.NewChain(nil, chainName)
	expectedDecoded := notSolo.ChainManager.MustGetContractID(chainName, accounts.Interface.Name)
	dataBytes := notSolo.DataManager.MustGetBytes(expectedDecoded)

	// Act
	actualDecoded := notSolo.DataManager.MustGetContractID(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}

func Test_MustGetHashResult(t *testing.T) {
	// Arrange
	notSolo := notsolo.New(t)
	expectedDecoded := notSolo.ChainManager.NewChain(nil, "dummyChain").State.Hash()
	dataBytes := notSolo.DataManager.MustGetBytes(&expectedDecoded)

	// Act
	actualDecoded := notSolo.DataManager.MustGetHash(dataBytes)

	// Assert
	require.Equal(t, expectedDecoded, actualDecoded)
}
