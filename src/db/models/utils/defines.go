package models_utils

type TransactionEvent string

const (
	CREATE TransactionEvent = "CREATE"
	UPDATE TransactionEvent = "UPDATE"
	DELETE TransactionEvent = "DELETE"
)

type MinifiedUser struct {
    Email    *string
    Password *string    
}
