package datamanager

import (
	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/iscp"
	"github.com/iotaledger/wasp/packages/iscp/colored"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/require"
)

const (
	couldNotConvertDataInto = "Could not convert data into "
	dataDoesNotExist        = "Data does not exist."
)

// DataManager manipulates result structures
type DataManager struct {
	env *solo.Solo
}

// New instantiates a data manager
func New(env *solo.Solo) *DataManager {
	resultHandler := &DataManager{env: env}
	return resultHandler
}

// MustGetInt64 converts input data into int64. Fails test when either no data is provided or cannot be converted.
func (dataManager *DataManager) MustGetInt64(data []byte) int64 {
	result, exists, err := codec.DecodeInt64(data)
	require.NoError(dataManager.env.T, err, couldNotConvertDataInto+" int64")
	require.True(dataManager.env.T, exists, dataDoesNotExist)
	return result
}

// MustGetString converts input data into int64. Fails test when either no data is provided or cannot be converted.
func (dataManager *DataManager) MustGetString(data []byte) string {
	result, exists, err := codec.DecodeString(data)
	require.NoError(dataManager.env.T, err, couldNotConvertDataInto+" string")
	require.True(dataManager.env.T, exists, dataDoesNotExist)
	return result
}

// MustGetAgentID converts input data into an AgentID. Fails test if no input is provided or cannot be converted.
func (dataManager *DataManager) MustGetAgentID(data []byte) iscp.AgentID {
	result, exists, err := codec.DecodeAgentID(data)
	require.NoError(dataManager.env.T, err, couldNotConvertDataInto+" AgentID")
	require.True(dataManager.env.T, exists, dataDoesNotExist)
	return result
}

// MustGetAddress converts input data into an Address. Fails test if no input is provided or cannot be converted.
func (dataManager *DataManager) MustGetAddress(data []byte) ledgerstate.Address {
	result, exists, err := codec.DecodeAddress(data)
	require.NoError(dataManager.env.T, err, couldNotConvertDataInto+" Address")
	require.True(dataManager.env.T, exists, dataDoesNotExist)
	return result
}

// MustGetChainID converts input data into a ChainID. Fails test if no input is provided or cannot be converted.
func (dataManager *DataManager) MustGetChainID(data []byte) iscp.ChainID {
	result, exists, err := codec.DecodeChainID(data)
	require.NoError(dataManager.env.T, err, couldNotConvertDataInto+" ChainID")
	require.True(dataManager.env.T, exists, dataDoesNotExist)
	return result
}

// MustGetColor converts input data into a Color. Fails test if no input is provided or cannot be converted.
func (dataManager *DataManager) MustGetColor(data []byte) colored.Color {
	result, exists, err := codec.DecodeColor(data)
	require.NoError(dataManager.env.T, err, couldNotConvertDataInto+" Color")
	require.True(dataManager.env.T, exists, dataDoesNotExist)
	return result
}

// MustGetHash converts input data into a HashValue. Fails test if no input is provided or cannot be converted.
func (dataManager *DataManager) MustGetHash(data []byte) hashing.HashValue {
	result, exists, err := codec.DecodeHashValue(data)
	require.NoError(dataManager.env.T, err, couldNotConvertDataInto+" HashValue")
	require.True(dataManager.env.T, exists, dataDoesNotExist)
	return result
}

// MustGetHname converts input data into an Hname. Fails test if no input is provided or cannot be converted.
func (dataManager *DataManager) MustGetHname(data []byte) iscp.Hname {
	result, exists, err := codec.DecodeHname(data)
	require.NoError(dataManager.env.T, err, couldNotConvertDataInto+" Hname")
	require.True(dataManager.env.T, exists, dataDoesNotExist)
	return result
}

// MustGetBytes returns the input as is. Fails test if no input is provided.
func (dataManager *DataManager) MustGetBytes(data interface{}) []byte {
	var bytes []byte
	require.NotPanics(dataManager.env.T, func() { bytes = codec.Encode(data) }, couldNotConvertDataInto+" bytes")
	require.NotNil(dataManager.env.T, bytes, couldNotConvertDataInto+" bytes")
	return bytes
}

// MustGetBool converts input data into a bool. Fails test if no input is provided or the data array is longer than 1.
func (dataManager *DataManager) MustGetBool(data []byte) bool {
	bytes := dataManager.MustGetBytes(data)
	require.Len(dataManager.env.T, bytes, 1, "Data does not have length 1")
	return bytes[0] != 0
}
