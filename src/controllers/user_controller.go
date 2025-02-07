package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alexbsec/MiniMarketplace/src/db/config"
	"github.com/alexbsec/MiniMarketplace/src/db/models"
	"golang.org/x/crypto/bcrypt"
)

var userService *models.UserService

// This is used to validate password 
type userBody struct {
    Name            *string `json:"name"`
    Email           *string `json:"email"`
    Password        *string `json:"password"`
    ConfirmPassword *string `json:"confirm_password"`
}

type userOut struct {
    ID          uint    `json:"id"`
    Name        *string `json:"name"`
    Email       *string `json:"email"`
}

func init() {
    service, err := config.InitService()
    if err != nil {
        panic(fmt.Sprintf("Failed to initialize database service: %v", err))
    }
    userService = &models.UserService{Service: service}
}

// Handles
func HandleUsers(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        handleCreateUser(w, r)
    case http.MethodGet:
        handleFetchUser(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
        
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
    var userIn userBody 
    if err := json.NewDecoder(r.Body).Decode(&userIn); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if userIn.Name == nil || userIn.Email == nil || userIn.Password == nil || userIn.ConfirmPassword == nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    hash, err := validateAndHashPassword(*userIn.Password, *userIn.ConfirmPassword)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to validate passwords: %v", err), http.StatusBadRequest)
        return
    }

    user := &models.User{
        Name:  userIn.Name, 
        Email:  userIn.Email,
        Password: hash,
    }

    user.Name = userIn.Name
    user.Email = userIn.Email
    user.Password = hash

    if err := userService.Create(user); err != nil {
        http.Error(w, "Failed to create user", http.StatusInternalServerError)
        return
    }

    var out userOut
    out.ID = user.ID
    out.Name = user.Name
    out.Email = user.Email

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(out)
}

func handleFetchUser(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid user", http.StatusBadRequest) 
        return
    }

    user, err := userService.Fetch(uint(id))
    if err != nil {
        http.Error(w, "User does not exist", http.StatusNotFound)
        return
    }

    if user == nil {
        http.Error(w, "User does not exist", http.StatusNotFound)
        return
    }

    var out userOut
    out.ID = user.ID
    out.Name = user.Name
    out.Email = user.Email

    json.NewEncoder(w).Encode(out)
}

func validateAndHashPassword(password string, confirmPassword string) (*string, error) {
    if password != confirmPassword {
        return nil, fmt.Errorf("passwords do not match")
    }
   
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    out := new(string)
    *out = string(bytes)
    return out, err
}

func verifyPassword(password string, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err != nil
}
