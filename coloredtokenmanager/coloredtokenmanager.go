package coloredtokenmanager

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
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

// MintColoredTokens converts a specified amount of balance of iota tokens available to SignatureScheme into a new color. Returns error if it fails.
func (coloredTokenmanager *ColoredTokenManager) MintColoredTokens(sigScheme signaturescheme.SignatureScheme, amount int64) (balance.Color, error) {
	return coloredTokenmanager.env.MintTokens(sigScheme, amount)
}

// MustMintColoredTokens converts a specified amount of balance of iota tokens available to SignatureScheme into a new color. Fails test on error.
func (coloredTokenmanager *ColoredTokenManager) MustMintColoredTokens(sigScheme signaturescheme.SignatureScheme, amount int64) balance.Color {
	color, err := coloredTokenmanager.MintColoredTokens(sigScheme, amount)
	require.NoError(coloredTokenmanager.env.T, err)
	return color
}

// DestroyColoredTokens converts a specified amount of balance of a specified color available to SignatureScheme into iota tokens. Returns error if it fails.
func (coloredTokenmanager *ColoredTokenManager) DestroyColoredTokens(sigScheme signaturescheme.SignatureScheme, color balance.Color, amount int64) error {
	return coloredTokenmanager.env.DestroyColoredTokens(sigScheme, color, amount)
}

// MustDestroyColoredTokens converts a specified amount of balance of a specified color available to SignatureScheme into iota tokens. Fails test on error.
func (coloredTokenmanager *ColoredTokenManager) MustDestroyColoredTokens(sigScheme signaturescheme.SignatureScheme, color balance.Color, amount int64) {
	err := coloredTokenmanager.env.DestroyColoredTokens(sigScheme, color, amount)
	require.NoError(coloredTokenmanager.env.T, err)
}
