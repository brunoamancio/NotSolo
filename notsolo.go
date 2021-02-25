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

// New instantiates NotSolo with default settings
func New(t *testing.T) *NotSolo {
	notSolo := &NotSolo{t: t}
	env := solo.New(t, notSolo.debug, notSolo.printStackTrace)

	notSolo.env = env
	notSolo.SigScheme = signatureschememanager.New(t, env)
	notSolo.ColoredToken = coloredtokenmanager.New(t, env)
	notSolo.Chain = chainmanager.New(t, env)
	notSolo.Request = requestmanager.New(t)
	notSolo.Response = responsemanager.New(t)
	notSolo.Data = datamanager.New(t)
	return notSolo
}
