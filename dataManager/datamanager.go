package datamanager

import (
	"testing"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/stretchr/testify/require"
)

const (
	couldNotConvertDataInto = "Could not convert data into "
	dataDoesNotExist        = "Data does not exist."
)

// DataManager manipulates result structures
type DataManager struct {
	t *testing.T
}

// New instantiates a data manager
func New(t *testing.T) *DataManager {
	resultHandler := &DataManager{t: t}
	return resultHandler
}

// MustGetInt64 converts input data into int64. Panics when either no data is provided or cannot be converted.
func (dataManager *DataManager) MustGetInt64(data []byte) int64 {
	result, exists, err := codec.DecodeInt64(data)
	require.NoError(dataManager.t, err, couldNotConvertDataInto+" int64")
	require.True(dataManager.t, exists, dataDoesNotExist)
	return result
}

// MustGetString converts input data into int64. Panics when either no data is provided or cannot be converted.
func (dataManager *DataManager) MustGetString(data []byte) string {
	result, exists, err := codec.DecodeString(data)
	require.NoError(dataManager.t, err, couldNotConvertDataInto+" string")
	require.True(dataManager.t, exists, dataDoesNotExist)
	return result
}

// MustGetAgentID converts input data into an AgentID. Panics if no input is provided or cannot be converted.
func (dataManager *DataManager) MustGetAgentID(data []byte) coretypes.AgentID {
	result, exists, err := codec.DecodeAgentID(data)
	require.NoError(dataManager.t, err, couldNotConvertDataInto+" AgentID")
	require.True(dataManager.t, exists, dataDoesNotExist)
	return result
}

// MustGetAddress converts input data into an Address. Panics if no input is provided or cannot be converted.
func (dataManager *DataManager) MustGetAddress(data []byte) address.Address {
	result, exists, err := codec.DecodeAddress(data)
	require.NoError(dataManager.t, err, couldNotConvertDataInto+" Address")
	require.True(dataManager.t, exists, dataDoesNotExist)
	return result
}

// MustGetChainID converts input data into a ChainID. Panics if no input is provided or cannot be converted.
func (dataManager *DataManager) MustGetChainID(data []byte) coretypes.ChainID {
	result, exists, err := codec.DecodeChainID(data)
	require.NoError(dataManager.t, err, couldNotConvertDataInto+" ChainID")
	require.True(dataManager.t, exists, dataDoesNotExist)
	return result
}

// MustGetColor converts input data into a Color. Panics if no input is provided or cannot be converted.
func (dataManager *DataManager) MustGetColor(data []byte) balance.Color {
	result, exists, err := codec.DecodeColor(data)
	require.NoError(dataManager.t, err, couldNotConvertDataInto+" Color")
	require.True(dataManager.t, exists, dataDoesNotExist)
	return result
}

// MustGetContractID converts input data into a ContractID. Panics if no input is provided or cannot be converted.
func (dataManager *DataManager) MustGetContractID(data []byte) coretypes.ContractID {
	result, exists, err := codec.DecodeContractID(data)
	require.NoError(dataManager.t, err, couldNotConvertDataInto+" ContractID")
	require.True(dataManager.t, exists, dataDoesNotExist)
	return result
}

// MustGetHash converts input data into a HashValue. Panics if no input is provided or cannot be converted.
func (dataManager *DataManager) MustGetHash(data []byte) hashing.HashValue {
	result, exists, err := codec.DecodeHashValue(data)
	require.NoError(dataManager.t, err, couldNotConvertDataInto+" HashValue")
	require.True(dataManager.t, exists, dataDoesNotExist)
	return *result
}

// MustGetHname converts input data into an Hname. Panics if no input is provided or cannot be converted.
func (dataManager *DataManager) MustGetHname(data []byte) coretypes.Hname {
	result, exists, err := codec.DecodeHname(data)
	require.NoError(dataManager.t, err, couldNotConvertDataInto+" Hname")
	require.True(dataManager.t, exists, dataDoesNotExist)
	return result
}

// MustGetBytes returns the input as is. Panics if no input is provided.
func (dataManager *DataManager) MustGetBytes(data interface{}) []byte {
	var bytes []byte
	require.NotPanics(dataManager.t, func() { bytes = codec.Encode(data) }, couldNotConvertDataInto+" bytes")
	require.NotNil(dataManager.t, bytes, couldNotConvertDataInto+" bytes")
	return bytes
}
