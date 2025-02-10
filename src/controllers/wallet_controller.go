package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexbsec/MiniMarketplace/src/db/models"
)

type walletOut struct {
    ID      uint
    Name    *string
    Amount  *float64
    Points  *float64
    UserID  uint
}

func HandleWallets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleCreateWallet(w, r)
    case http.MethodGet:
        handleFetchWallet(w, r)
    case http.MethodPut:
        handleUpdateWallet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleCreateWallet(w http.ResponseWriter, r *http.Request) {
    var wallet models.Wallet
    if err := json.NewDecoder(r.Body).Decode(&wallet); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // These fields must not be empty
    if wallet.Name == nil || wallet.Amount == nil || wallet.Points == nil {
        http.Error(w, "Invalid request paylaod", http.StatusBadRequest)
        return
    }

    // Authentication step
    user, result := UserAuthFlowLax(w, r, ROLE_USER)
    if !result {
        return
    }

    // If user is not admin, the parameter 'user_id' must
    // not be parsed
    if *user.Role == uint(ROLE_ADMIN) {
        // Just assigns the User correctly
        assignedUser, err := userService.Fetch(wallet.UserID)
        if err != nil {
            http.Error(w, "User does not exist", http.StatusBadRequest)
            return
        }
        wallet.User = *assignedUser
    } else {
        // Assigns the fetched userID to this wallet
        wallet.UserID = user.ID
        wallet.User = *user

        // Ensure no amount is parsed for the creation of a new wallet
        if *wallet.Amount != 0  {
            *wallet.Amount = 0
        }

        if *wallet.Points != 0 {
            *wallet.Points = 0
        } 
    }

    if err := walletService.Create(&wallet); err != nil {
        http.Error(w, "Failed to create wallet", http.StatusInternalServerError)
        return
    }

    var out walletOut
    out.ID = wallet.ID
    out.Name = wallet.Name
    out.Points = wallet.Points
    out.Amount = wallet.Amount
    out.UserID = wallet.UserID

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(out)
}

func handleFetchWallet(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid wallet", http.StatusBadRequest) 
        return
    }

    user, result := UserAuthFlowLax(w, r, ROLE_USER)
    if !result {
        return
    }

    wallet, err := walletService.Fetch(uint(id))
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    if *user.Role != uint(ROLE_ADMIN) && wallet.UserID != user.ID {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var out walletOut
    out.ID = wallet.ID
    out.Name = wallet.Name
    out.Amount = wallet.Amount
    out.Points = wallet.Points
    out.UserID = wallet.UserID

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(out)
}

func handleUpdateWallet(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid wallet", http.StatusBadRequest) 
        return
    }

    user, result := UserAuthFlowLax(w, r, ROLE_USER)
    if !result {
        return
    }

    var newWallet models.Wallet
    if err := json.NewDecoder(r.Body).Decode(&newWallet); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    wallet, err := walletService.Fetch(uint(id))
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    if newWallet.UserID != wallet.UserID {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    if *user.Role != uint(ROLE_ADMIN) && wallet.UserID != user.ID {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Only allow amount and points update if request is made by admin
    if *user.Role == uint(ROLE_ADMIN) {
        if newWallet.Amount != nil {
           wallet.Amount = newWallet.Amount 
        }

        if newWallet.Points != nil {
            wallet.Points = newWallet.Points
        }
    }

    if newWallet.Name != nil {
        wallet.Name = newWallet.Name
    }


    if err := walletService.Update(uint(id), wallet); err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var out walletOut
    out.ID = wallet.ID
    out.Name = wallet.Name
    out.Amount = wallet.Amount
    out.Points = wallet.Points
    out.UserID = wallet.UserID

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(out)
}
