package notsolo

import (
	"testing"

	chainmanager "github.com/brunoamancio/NotSolo/chainManager"
	coloredtokenmanager "github.com/brunoamancio/NotSolo/coloredTokenManager"
	datamanager "github.com/brunoamancio/NotSolo/dataManager"
	requestmanager "github.com/brunoamancio/NotSolo/requestManager"
	responsemanager "github.com/brunoamancio/NotSolo/responseManager"
	signatureschememanager "github.com/brunoamancio/NotSolo/signatureSchemeManager"
	"github.com/iotaledger/wasp/packages/solo"
)

// NotSolo is a wrapper around solo to simplify unit testing
type NotSolo struct {
	t               *testing.T
	debug           bool
	printStackTrace bool
	env             *solo.Solo
	SigScheme       *signatureschememanager.SignatureSchemeManager
	ColoredToken    *coloredtokenmanager.ColoredTokenManager
	Chain           *chainmanager.ChainManager
	Request         *requestmanager.RequestManager
	Response        *responsemanager.ResponseManager
	Data            *datamanager.DataManager
}

var notSolo *NotSolo = &NotSolo{}

// Initializable defines a contract to verify whether a structure is initialized
type Initializable interface {
	IsInitialized() bool
}

// IsInitialized implements Initializable for *NotSolo
func (notSolo *NotSolo) IsInitialized() bool {
	return notSolo.env != nil
}

// New instantiates NotSolo with default settings
func New(t *testing.T) *NotSolo {
	if notSolo.IsInitialized() {
		notSolo.env.T = t
	} else {
		loadManagers(t)
	}

	notSolo.t = t
	return notSolo
}

func loadManagers(t *testing.T) {
	notSolo.env = solo.New(t, notSolo.debug, notSolo.printStackTrace)

	notSolo.SigScheme = signatureschememanager.New(notSolo.env)
	notSolo.ColoredToken = coloredtokenmanager.New(notSolo.env)
	notSolo.Chain = chainmanager.New(notSolo.env)
	notSolo.Request = requestmanager.New(notSolo.env)
	notSolo.Response = responsemanager.New(notSolo.env)
	notSolo.Data = datamanager.New(notSolo.env)
}
