package urlx

import (
	"encoding/base64"
)

/* headers */

func Authorization(authorization string) Option {
	if authorization != "" {
		return HeaderSet("Authorization", authorization)
	}
	return HeaderSet("-Authorization", "")
}

func Bearer(bearerToken string) Option {
	if bearerToken != "" {
		Authorization("Bearer " + bearerToken)
	}
	return HeaderSet("-Authorization", "")
}

func BasicAuth(user, pass string) Option {
	if pass != "" {
		return Authorization(base64.StdEncoding.EncodeToString([]byte(user + ":" + pass)))
	}
	return HeaderSet("-Authorization", "")
}
