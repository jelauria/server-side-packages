package handlers

import (
	"assignments-jelauria/servers/gateway/models/users"
	"assignments-jelauria/servers/gateway/sessions"
)

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store
type Context struct {
	SeshKey   string
	SeshStore *sessions.Store
	UserStore *users.Store
}
