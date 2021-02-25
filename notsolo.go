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
	t                      *testing.T
	debug                  bool
	printStackTrace        bool
	env                    *solo.Solo
	SignatureSchemeManager *signatureschememanager.SignatureSchemeManager
	ColoredTokenManager    *coloredtokenmanager.ColoredTokenManager
	ChainManager           *chainmanager.ChainManager
	RequestManager         *requestmanager.RequestManager
	ResponseManager        *responsemanager.ResponseManager
	DataManager            *datamanager.DataManager
}

// New instantiates NotSolo with default settings
func New(t *testing.T) *NotSolo {
	notSolo := &NotSolo{t: t}
	env := solo.New(t, notSolo.debug, notSolo.printStackTrace)

	notSolo.env = env
	notSolo.SignatureSchemeManager = signatureschememanager.New(t, env)
	notSolo.ColoredTokenManager = coloredtokenmanager.New(t, env)
	notSolo.ChainManager = chainmanager.New(t, env)
	notSolo.RequestManager = requestmanager.New(t)
	notSolo.ResponseManager = responsemanager.New(t)
	notSolo.DataManager = datamanager.New(t)
	return notSolo
}
