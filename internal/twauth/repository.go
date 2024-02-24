package twauth

import "github.com/vctaragao/twitch-chat/internal/twauth/entity"

type Repository interface {
	InsertAuth(auth entity.AuthToken) (int, error)
	InsertScopes(auth entity.AuthToken) error
}
