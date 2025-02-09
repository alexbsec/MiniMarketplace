package controllers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/alexbsec/MiniMarketplace/src/db/models"
	"github.com/alexbsec/MiniMarketplace/src/logging"
	"github.com/golang-jwt/jwt/v5"
)

type Role uint

const (
    ROLE_ADMIN Role = 1
    ROLE_USER  Role = 0
)

// TODO: make this secure
var secretKey = []byte("test-secret")

type loginBody struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type JWTContent struct {
	ID       uint
	Email    string
	ExpireAt float64
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleLoginUser(w, r)
	default:
		http.Error(w, "Method not allowd", http.StatusMethodNotAllowed)
	}
}

func UserAuthFlowLax(w http.ResponseWriter, r *http.Request, expectRole Role) (*models.User, bool) {
    jwtCtt, err := parseJWT(r) 
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return nil, false 
    }

    reqUser, err := userService.Fetch(jwtCtt.ID) 
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return nil, false
    }
  
    _, err = validateUserRequest(r, reqUser, expectRole) 
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return nil, false
    }

    return reqUser, true
}

func UserAuthFlow(w http.ResponseWriter, r *http.Request, id uint, expectRole Role) bool {
    jwtCtt, err := parseJWT(r) 
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return false
    }

    if jwtCtt.ID != id {
        http.Error(w, "Unauthorized", http.StatusUnauthorized) 
        return false
    } 

    reqUser, err := userService.Fetch(jwtCtt.ID) 
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return false
    }
  
    _, err = validateUserRequest(r, reqUser, expectRole) 
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return false
    }

    return true
}

func handleLoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
	if !passwordsMatch {
		http.Error(w, "Usuário ou senha incorretos", http.StatusUnauthorized)
		return
	}

	// TODO: Handle login logic (give token etc)
	tokenString, err := createJWT(userRec)
	if err != nil {
		http.Error(w, "Failed to create session token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func validateUserRequest(r *http.Request, user *models.User, expectRole Role) (*JWTContent, error) {
    if *user.Role < uint(expectRole) {
        return nil, fmt.Errorf("Unauthorized") 
    }

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("Unauthorized")
	}

    tokenString = tokenString[len("Bearer "):]
	token, err := verifyJWT(user, tokenString)
	if err != nil {
		logging.Log.Warn("Verification of JWT failed", slog.String("cause", err.Error()))
		return nil, fmt.Errorf("Unauthorized")
	}

	return token, nil
}

func parseJWT(r *http.Request) (*JWTContent, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("Unauthorized")
	}

    tokenString = tokenString[len("Bearer "):]
    logging.Log.Debug("JWT token being processed", slog.String("token", tokenString))

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })

    if err != nil {
        logging.Log.Error("Failed to verify JWT token", slog.String("error", err.Error()))
        return nil, fmt.Errorf("Unauthorized")
    }

    var jwtContent JWTContent
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid token expiration")
		}
        jwtContent.ExpireAt = exp

		// Extract and validate "email" claim
		email, ok := claims["email"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid user email")
		}
        jwtContent.Email = email

		// Extract and validate "id" claim
		idFloat, ok := claims["id"].(float64)
		if !ok {
            logging.Log.Debug("Invalid user ID", slog.String("id", strconv.Itoa(int(idFloat))))
			return nil, fmt.Errorf("invalid user ID")
		}
        jwtContent.ID = uint(idFloat)
    } else {
		return nil, fmt.Errorf("invalid token claims")
	}

    return &jwtContent, nil
}

func createJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":    user.ID,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
			"email": *user.Email,
		},
	)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		logging.Log.Error("Failed to create JWT token", slog.String("error", err.Error()))
		return "", err
	}

	return tokenString, nil
}


func verifyJWT(user *models.User, tokenString string) (*JWTContent, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })

    if err != nil {
        logging.Log.Error("Failed to verify JWT token", slog.String("error", err.Error()))
        return nil, err
    }

    var jwtContent JWTContent
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, fmt.Errorf("invalid token expiration")
		}
        jwtContent.ExpireAt = exp

		// Extract and validate "email" claim
		email, ok := claims["email"].(string)
		if !ok || user.Email == nil || email != *user.Email {
			return nil, fmt.Errorf("invalid user email")
		}
        jwtContent.Email = email

		// Extract and validate "id" claim
		idFloat, ok := claims["id"].(float64)
		if !ok || uint(idFloat) != user.ID {
			return nil, fmt.Errorf("invalid user ID")
		}
        jwtContent.ID = uint(idFloat)

		// Check if token is expired
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, fmt.Errorf("token has expired")
		}
	} else {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &jwtContent, nil
}
