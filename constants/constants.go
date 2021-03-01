package constants

const (
	// IotaTokensConsumedByRequest IMPORTANT: When a request is sent to a chain, 1 IOTA is sent from the caller's address in the value tangle to their account in the chain
	IotaTokensConsumedByRequest = 1
	// IotaTokensConsumedByChain IMPORTANT: When a chain is created, 1 IOTA is colored with the chain's color and sent to the chain's address in the value tangle
	IotaTokensConsumedByChain = 1
)
