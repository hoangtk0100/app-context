package core

import "context"

const KeyRequester = "requester"

type Requester interface {
	GetID() string
	GetUID() string
}

type requesterData struct {
	ID  string `json:"id"`
	UID string `json:"uid"`
}

func NewRequester(id, uid string) *requesterData {
	return &requesterData{
		ID:  id,
		UID: uid,
	}
}

func (r *requesterData) GetID() string {
	return r.ID
}

func (r *requesterData) GetUID() string {
	return r.UID
}

func GetRequester(ctx context.Context) Requester {
	if requester, ok := ctx.Value(KeyRequester).(Requester); ok {
		return requester
	}

	return nil
}

func ContextWithRequester(ctx context.Context, requester Requester) context.Context {
	return context.WithValue(ctx, KeyRequester, requester)
}
