package notsolo

import (
	"testing"

	"github.com/brunoamancio/NotSolo/chainmanager"
	"github.com/brunoamancio/NotSolo/coloredtokenmanager"
	"github.com/brunoamancio/NotSolo/datamanager"
	"github.com/brunoamancio/NotSolo/keypairmanager"
	"github.com/brunoamancio/NotSolo/l1manager"
	"github.com/brunoamancio/NotSolo/requestmanager"
	"github.com/iotaledger/wasp/packages/solo"
)

// NotSolo is a wrapper around solo to simplify unit testing
type NotSolo struct {
	t               *testing.T
	debug           bool
	printStackTrace bool
	env             *solo.Solo
	KeyPair         *keypairmanager.KeyPairManager
	ColoredToken    *coloredtokenmanager.ColoredTokenManager
	L1              *l1manager.L1Manager
	Chain           *chainmanager.ChainManager
	Request         *requestmanager.RequestManager
	Data            *datamanager.DataManager
}

var notSolo *NotSolo = &NotSolo{}

// Initializable defines a contract to verify whether a structure is initialized
type Initializable interface {
	IsInitialized() bool
}

// Disposable defines a contract to clear resources
type Disposable interface {
	Dispose()
}

// IsInitialized implements Initializable for *NotSolo
func (notSolo *NotSolo) IsInitialized() bool {
	return notSolo.env != nil
}

// Dispose implements Disposable for NotSolo
func (notSolo *NotSolo) Dispose() {
	notSolo.Chain.Dispose()
}

// New instantiates NotSolo with default settings
func New(t *testing.T) *NotSolo {
	if notSolo.IsInitialized() {
		notSolo.Dispose()
		notSolo.env.T = t
	} else {
		loadManagers(t)
	}

	notSolo.t = t
	return notSolo
}

func loadManagers(t *testing.T) {
	notSolo.env = solo.New(t, notSolo.debug, notSolo.printStackTrace)

	notSolo.KeyPair = keypairmanager.New(notSolo.env)
	notSolo.ColoredToken = coloredtokenmanager.New(notSolo.env)
	notSolo.Chain = chainmanager.New(notSolo.env)
	notSolo.L1 = l1manager.New(notSolo.env)
	notSolo.Request = requestmanager.New(notSolo.env)
	notSolo.Data = datamanager.New(notSolo.env)
}
