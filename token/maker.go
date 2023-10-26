package token

import (
	"time"
)

// jwt maker と paseto maker のインターフェースを定義
type Maker interface {
	// create a token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// check if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
