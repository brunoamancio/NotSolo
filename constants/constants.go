package constants

const (
	// DefaultChainStartingBalance IMPORTANT: When a chain is created >>> using SOLO<<<, the chainOriginator sends 100 tokens to chainID
	DefaultChainStartingBalance = uint64(100)

	// IotaTokensConsumedByRequest IMPORTANT: When a request is sent to a chain, 1 IOTA is sent from the caller's address in L1 to their account in the chain
	IotaTokensConsumedByRequest = 1
)
