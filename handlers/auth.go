package handlers

import (
	"assignments-jelauria/servers/gateway/models/users"
	"assignments-jelauria/servers/gateway/sessions"
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.
func (c *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		conentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, errors.New("Request body must be in JSON."), http.StatusUnsupportedMediaType)
			return
		}
		var nu users.NewUser
		jsonErr := json.NewDecoder(r.Body).Decode(&nu)
		if err != nil {
			http.Error(w, jsonErr.Error(), http.StatusBadRequest)
			return
		}
		user, toUsrErr := nu.ToUser()
		if toUsrErr != nil {
			http.Error(w, toUsrErr.Error(), http.StatusBadRequest)
			return
		}
		authUsr, insertErr := c.UserStore.Insert(user)
		if insertErr != nil {
			http.Error(w, insertErr.Error(), http.StatusBadRequest)
			return
		}
		_, keyErr := sessions.BeginSession(c.SeshKey, c.SeshStore, &SessionState{time.Now(), authUsr}, w)
		if keyErr != nil {
			http.Error(w, keyErr.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(authUsr)
		return
	}
	http.Error(w, errors.New("Http method not allowed."), http.StatusMethodNotAllowed)
	return
}

func (c *Context) SpecificUsersHandler(w http.ResponseWriter, r *http.Request) {
	seshID, seshErr := sessions.GetSessionID(r, c.SeshKey)
	if seshErr != nil {
		http.Error(w, seshErr.Error(), http.StatusUnauthorized)
		return
	}
	currState := &SessionState{}
	valErr := c.SeshStore.Get(seshID, currState)
	if valErr != nil {
		http.Error(w, valErr.Error(), http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodGet {
		strID := path.Base(r.URL.Path)
		if strID == "me" {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(currState.AuthUser)
			return
		}
		intID, strErr := strconv.ParseInt(strID, 10, 64)
		if strErr != nil {
			http.Error(w, errors.New("No user with given ID."), http.StatusNotFound)
		}
		qUser, sqlErr := c.UserStore.GetByID(intID)
		if sqlErr != nil {
			http.Error(w, sqlErr.Error(), http.StatusNotFound)
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(qUser)
		return
	}
	if r.Method == http.MethodPatch {
		pID := path.Base(r.URL.Path)
		if pID != "me" {
			http.Error(w, errors.New("Forbidden request."), http.StatusForbidden)
		}
		conentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, errors.New("Request body must be in JSON."), http.StatusUnsupportedMediaType)
			return
		}
		var updates users.Updates
		jsonErr := json.NewDecoder(r.Body).Decode(&updates)
		if err != nil {
			http.Error(w, jsonErr.Error(), http.StatusBadRequest)
			return
		}
		updUser, upErr := c.UserStore.Update(currState.AuthUser.ID, updates)
		if upErr != nil {
			http.Error(w, upErr.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updUser)
		return
	}
	http.Error(w, errors.New("Http method not allowed."), http.StatusMethodNotAllowed)
	return
}

func (c *Context) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		conentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, errors.New("Request body must be in JSON."), http.StatusUnsupportedMediaType)
			return
		}
		var creds *users.Credentials
		jsonErr := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, jsonErr.Error(), http.StatusBadRequest)
			return
		}
		user, getErr := c.UserStore.GetByEmail(creds.Email)
		if getErr != nil {
			fakeUsr := &users.User{-123, "helloworld", []byte("loveUAvatarAang"), "Toph", "Beifong", "blah"}
			fakeUsr.Authenticate("12345")
			http.Error(w, getErr.Error(), http.StatusUnauthorized)
			return
		}
		authErr := user.Authenticate(creds.Password)
		if authErr != nil {
			http.Error(w, getErr.Error(), http.StatusUnauthorized)
			return
		}
		_, keyErr := sessions.BeginSession(c.SeshKey, c.SeshStore, &SessionState{time.Now(), user}, w)
		if keyErr != nil {
			http.Error(w, keyErr.Error(), http.StatusUnauthorized)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
		return
	}
	http.Error(w, errors.New("Http method not allowed."), http.StatusMethodNotAllowed)
	return
}

func (c *Context) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		pID := path.Base(r.URL.Path)
		if pID != "mine" {
			http.Error(w, errors.New("Forbidden request."), http.StatusForbidden)
		}
		sessions.EndSession(r, c.SeshKey, c.SeshStore)
		w.Write([]byte("signed out"))
		return
	}
	http.Error(w, errors.New("Http method not allowed."), http.StatusMethodNotAllowed)
	return
}
