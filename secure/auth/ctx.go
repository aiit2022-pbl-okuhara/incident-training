package auth

import (
	"context"

	"github.com/google/go-safeweb/safehttp"
)

type key string

const (
	userKey       key = "user"
	changeSessKey key = "change"
)

type sessionAction string

const (
	clearSess sessionAction = "clear"
	setSess   sessionAction = "set"
)

func putSessionAction(ctx context.Context, action sessionAction) {
	safehttp.FlightValues(ctx).Put(changeSessKey, action)
}

func ctxSessionAction(ctx context.Context) sessionAction {
	v := safehttp.FlightValues(ctx).Get(changeSessKey)
	action, ok := v.(sessionAction)
	if !ok {
		return ""
	}
	return action
}

func putUser(ctx context.Context, user string) {
	safehttp.FlightValues(ctx).Put(userKey, user)
}

func ctxUser(ctx context.Context) string {
	v := safehttp.FlightValues(ctx).Get(userKey)
	user, ok := v.(string)
	if !ok {
		return ""
	}
	return user
}
