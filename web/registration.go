package web

import ( 
	"net/http"
	"fmt"
	"encoding/json" 
	"encoding/hex"

	"github.com/SCKelemen/Cassius/data"
)

func RegisterHandler(w http.ResponseWriter, req *http.Request, env *environment) {
	var registration struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&registration); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Error decoding request: %v", err)
		return
	}

	if registration.Name == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintln(w, `Request must include the attribute "name"`)
		return
	}

	if len(registration.Name) > 30 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintln(w, `"name" must be less than 30 characters`)
		return
	}

	err := data.ValidatePassword(registration.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintln(w, err)
		return
	}

	user := data.User{}
	user.Name = data.NewString(registration.Name)
	user.Email = data.NewStringFallback(registration.Email, data.Undefined)
	user.SetPassword(registration.Password)

	userID, err := data.CreateUser(env.pool, &user)
	if err != nil {
		if err, ok := err.(data.DuplicationError); ok {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, `"%s" is already taken`, err.Field)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var response struct {
		Name   string `json:"name"`
		UserID int32  `json:"userID"`
	}

	response.Name = registration.Name
	response.UserID = userID

	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

func DeregisterHandler(w http.ResponseWriter, req *http.Request, env *environment) {
	sessionID, err := hex.DecodeString(req.Header.Get("X-Authentication"))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = data.DeleteSession(env.pool, sessionID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = data.DeleteUser(env.pool, env.user.ID.Value)
	if err != nil {
		http.Error(w, "Unable to find user", http.StatusInternalServerError)
		return
	}
}
