package coloredtokenmanager

import (
	"github.com/iotaledger/hive.go/crypto/ed25519"
	"github.com/iotaledger/wasp/packages/iscp/colored"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/stretchr/testify/require"
)

// ColoredTokenManager manipulates colored tokens
type ColoredTokenManager struct {
	env *solo.Solo
}

// New instantiates a colored token manager
func New(env *solo.Solo) *ColoredTokenManager {
	coloredTokenManager := &ColoredTokenManager{env: env}
	return coloredTokenManager
}

// MintColoredTokens converts a specified amount of balance of iota tokens available to ed25519.KeyPair into a new color. Returns error if it fails.
func (coloredTokenmanager *ColoredTokenManager) MintColoredTokens(keyPair *ed25519.KeyPair, amount uint64) (colored.Color, error) {
	return coloredTokenmanager.env.MintTokens(keyPair, amount)
}

// MustMintColoredTokens converts a specified amount of balance of iota tokens available to ed25519.KeyPair into a new color. Fails test on error.
func (coloredTokenmanager *ColoredTokenManager) MustMintColoredTokens(keyPair *ed25519.KeyPair, amount uint64) colored.Color {
	color, err := coloredTokenmanager.MintColoredTokens(keyPair, amount)
	require.NoError(coloredTokenmanager.env.T, err)
	return color
}
