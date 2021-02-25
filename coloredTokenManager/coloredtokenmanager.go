package coloredtokenmanager

import (
	"testing"

	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/packages/solo"
)

// ColoredTokenManager manipulates colored tokens
type ColoredTokenManager struct {
	t   *testing.T
	env *solo.Solo
}

// New instantiates a colored token manager
func New(t *testing.T, env *solo.Solo) *ColoredTokenManager {
	coloredTokenManager := &ColoredTokenManager{t: t, env: env}
	return coloredTokenManager
}

// MintColoredTokens converts a specified amount of balance of iota tokens available to SignatureScheme into a new color. Returns error if it fails.
func (coloredTokenmanager *ColoredTokenManager) MintColoredTokens(sigScheme signaturescheme.SignatureScheme, amount int64) (balance.Color, error) {
	return coloredTokenmanager.env.MintTokens(sigScheme, amount)
}

// DestroyColoredTokens converts a specified amount of balance of a specified color available to SignatureScheme into iota tokens. Returns error if it fails.
func (coloredTokenmanager *ColoredTokenManager) DestroyColoredTokens(sigScheme signaturescheme.SignatureScheme, color balance.Color, amount int64) error {
	return coloredTokenmanager.env.DestroyColoredTokens(sigScheme, color, amount)
}
