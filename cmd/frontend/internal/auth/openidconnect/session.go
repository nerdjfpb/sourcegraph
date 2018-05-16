package openidconnect

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sourcegraph/sourcegraph/cmd/frontend/internal/session"
)

const sessionKey = "oidc@0"

type sessionData struct {
	Issuer   string
	ClientID string

	// Store only the oauth2.Token fields we need, to avoid hitting the ~4096-byte session data
	// limit.
	AccessToken string
	TokenType   string
}

func (d sessionData) toProviderKey() providerKey {
	return providerKey{issuerURL: d.Issuer, clientID: d.ClientID}
}

// SignOut clears OpenID Connect-related data from the session. If possible, it revokes the token
// from the OP. If there is an end-session endpoint, it returns that for the caller to present to
// the user.
func SignOut(w http.ResponseWriter, r *http.Request) (endSessionEndpoint string, err error) {
	defer func() {
		if err != nil {
			_ = session.SetData(w, r, sessionKey, nil) // clear the bad data
		}
	}()

	var data *sessionData
	if err := session.GetData(r, sessionKey, &data); err != nil {
		return "", errors.WithMessage(err, "reading OpenID Connect session data")
	}
	if err := session.SetData(w, r, sessionKey, nil); err != nil {
		return "", errors.WithMessage(err, "clearing OpenID Connect session data")
	}
	if data != nil {
		pc := getProviderConfig(data.toProviderKey())
		if pc == nil {
			return "", fmt.Errorf("unable to revoke token or end session for OpenID Connect because site config no longer contains issuer %q and client ID %q", data.Issuer, data.ClientID)
		}
		p, err := cache.get(data.Issuer)
		if err != nil {
			return "", errors.WithMessage(err, "looking up OpenID Connect provider metadata")
		}

		endSessionEndpoint = p.EndSessionEndpoint

		if p.RevocationEndpoint != "" {
			if err := revokeToken(r.Context(), pc, p.RevocationEndpoint, data.AccessToken, data.TokenType); err != nil {
				return endSessionEndpoint, errors.WithMessage(err, "revoking OpenID Connect token")
			}
		}
	}

	return endSessionEndpoint, nil
}
