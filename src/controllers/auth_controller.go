package controllers

import (
	"encoding/json"
	"net/http"
)

type loginBody struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleLoginUser(w, r)
	default:
		http.Error(w, "Method not allowd", http.StatusMethodNotAllowed)
	}
}

func handleLoginUser(w http.ResponseWriter, r *http.Request) {
    var params loginBody
    if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }
   
    if params.Email == nil || params.Password == nil {
        http.Error(w, "Please provide your e-mail and password", http.StatusBadRequest)
        return
    }

    userRec, err := userService.FetchUserByEmail(params.Email)
    if err != nil {
        http.Error(w, "Usuário ou senha incorretos", http.StatusNotFound)
        return
    }

    passwordsMatch := verifyPassword(*params.Password, *userRec.Password)
    if (!passwordsMatch) {
        http.Error(w, "Usuário ou senha incorretos", http.StatusUnauthorized) 
        return
    }

    // TODO: Handle login logic (give token etc)
    w.WriteHeader(http.StatusOK)
}
