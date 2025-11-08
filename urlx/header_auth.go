package urlx

import (
	"encoding/base64"
)

/* headers */

func Authorization(authorization string) HeaderOption {
	if authorization != "" {
		return HeaderSet("Authorization", authorization)
	} else {
		return HeaderDel("Authorization")
	}
}

func Bearer(bearerToken string) HeaderOption {
	if bearerToken != "" {
		Authorization("Bearer " + bearerToken)
	}
	return HeaderDel("Authorization")
}

func BasicAuth(user, pass string) HeaderOption {
	if user != "" || pass != "" {
		return Authorization(base64.StdEncoding.EncodeToString([]byte(user + ":" + pass)))
	}
	return HeaderDel("Authorization")
}
