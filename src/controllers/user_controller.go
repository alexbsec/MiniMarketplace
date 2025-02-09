package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alexbsec/MiniMarketplace/src/db/models"
	"golang.org/x/crypto/bcrypt"
)

// This is used to validate password
type userBody struct {
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	Password        *string `json:"password"`
	ConfirmPassword *string `json:"confirm_password"`
}

type userUpdateBody struct {
	Name            *string `json:"name"`
	Email           *string `json:"email"`
	OldPassword     *string `json:"old_password"`
	NewPassword     *string `json:"new_password"`
	ConfirmPassword *string `json:"confirm_password"`
}

type userOut struct {
	ID    uint    `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

// Handles
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleCreateUser(w, r)
	case http.MethodGet:
		handleFetchUser(w, r)
	case http.MethodPut:
		handleUpdateUser(w, r)
	case http.MethodDelete:
		handleDeleteUser(w, r)
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
		Name:     userIn.Name,
		Email:    userIn.Email,
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

func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusNotFound)
		return
	}

	user, err := userService.Fetch(uint(id))
	if err != nil {
		http.Error(w, "Usuário não existe", http.StatusNotFound)
		return
	}

	if user == nil {
		http.Error(w, "Usuário não existe", http.StatusNotFound)
		return
	}

	var updatedUser userUpdateBody
	if err = json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid parameters when updating", http.StatusBadRequest)
		return
	}

	result, status, err := checkUpdatePasswordFlow(uint(id), &updatedUser)
	if result != nil || status != http.StatusOK || err != nil {
		if err != nil && result == nil {
			// In this case we wont output the error message and say something instead
			http.Error(w, "Something unexpected happened", status)
			return
		} else if err != nil {
			// Here the result is a message but we got an error
			http.Error(w, *result, status)
			return
		}

		// It means we do not get any errors, so we can go and
		// update the password
		user.Password = result
	}

	if updatedUser.Name != nil {
		user.Name = updatedUser.Name
	}

	if updatedUser.Email != nil {
		result, status, err := checkUpdateEmailFlow(&updatedUser)
		if result == nil && status == http.StatusOK && err == nil {
			user.Email = updatedUser.Email
		} else {
			http.Error(w, *result, status)
			return
		}
	}

	if err = userService.Update(uint(id), user); err != nil {
		http.Error(w, "Falha ao atualizar usuário", http.StatusInternalServerError)
		return
	}

	out := &userOut{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(out)
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusNotFound)
		return
	}

	user, err := userService.Fetch(uint(id))
	if err != nil {
		http.Error(w, "Usuário não existe", http.StatusNotFound)
		return
	}

	if user == nil {
		http.Error(w, "Usuário não existe", http.StatusNotFound)
		return
	}

	if err = userService.Delete(uint(id)); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}



func hashPassword(password string) (*string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	out := new(string)
	*out = string(bytes)
	return out, err
}

func validateAndHashPassword(password string, confirmPassword string) (*string, error) {
	if password != confirmPassword {
		return nil, fmt.Errorf("passwords do not match")
	}

	return hashPassword(password)
}

func verifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func checkUpdatePasswordFlow(id uint, updateBody *userUpdateBody) (*string, int, error) {
	if updateBody == nil {
		errMsg := "Something unexpected happened"
		return &errMsg, http.StatusBadRequest, fmt.Errorf("updateBody is a nullptr")
	}

    // No update
    if updateBody.OldPassword == nil && updateBody.NewPassword == nil && updateBody.ConfirmPassword == nil {
        return nil, http.StatusOK, nil
    }

	// If old password is not provided but new/confirm password is, return error
	if updateBody.OldPassword == nil && (updateBody.NewPassword != nil || updateBody.ConfirmPassword != nil) {
		errMsg := "Old password is required to update to a new password"
		return &errMsg, http.StatusBadRequest, fmt.Errorf("Old password is required to update to a new password")
	}

	// If old password is provided but one of the new passwords is missing, return error
	if updateBody.OldPassword != nil && (updateBody.NewPassword == nil || updateBody.ConfirmPassword == nil) {
		errMsg := "Both 'new_password' and 'confirm_password' are required when updating a password"
		return &errMsg, http.StatusBadRequest, fmt.Errorf("Fields are required")
	}

	// Handle the password update flow
	if updateBody.OldPassword != nil && updateBody.NewPassword != nil && updateBody.ConfirmPassword != nil {
		// Fetch the user's current hashed password
		oldHash, err := userService.FetchPassword(id)
		if err != nil {
			errMsg := "Failed to fetch user password"
			return &errMsg, http.StatusInternalServerError, err
		}

		// Check if the old password matches
		if oldHash == nil || !verifyPassword(*updateBody.OldPassword, *oldHash) {
			errMsg := "Old password does not match our records"
			return &errMsg, http.StatusBadRequest, fmt.Errorf("Old password does not match our records")
		}

		// Validate and hash the new password
		newHash, err := validateAndHashPassword(*updateBody.NewPassword, *updateBody.ConfirmPassword)
		if err != nil {
			errMsg := "New passwords do not match"
			return &errMsg, http.StatusBadRequest, fmt.Errorf("New passwords do not match")
		}

		// Return the new password hash if everything is valid
		return newHash, http.StatusOK, nil
	}

	// This should never be reached if the input validation is correct
	errMsg := "Unexpected password update flow"
	return &errMsg, http.StatusInternalServerError, fmt.Errorf("unexpected password update flow")
}

func checkUpdateEmailFlow(updateBody *userUpdateBody) (*string, int, error) {
	if updateBody == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("updateBody is a nullptr")
	}

	if updateBody.Email == nil {
		return nil, http.StatusOK, nil
	}

	userExists, err := userService.CheckEmailExists(updateBody.Email)
	if err != nil {
		errMsg := "Something unexpected happen"
		return &errMsg, http.StatusInternalServerError, err
	}

	if userExists {
		errMsg := "Impossível alterar email pois email já existe"
		return &errMsg, http.StatusBadRequest, fmt.Errorf("Failed to update")
	}

	return nil, http.StatusOK, nil
}
