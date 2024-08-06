package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"account-service/data"
	"account-service/event"

	"github.com/go-chi/chi/v5"
)

const (
	DifferentBankFee = 4.5
)

// LogPayload is the embedded type (in RequestPayload) that describes a request to log something
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// MailPayload is the embedded type (in RequestPayload) that describes an email message to be sent
type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// TransferMoney handles the money transfer between accounts.
func (app *Config) TransferMoney(w http.ResponseWriter, r *http.Request) {
	sameBank := true
	var requestPayload struct {
		Sender   string  `json:"sender"`   // ID of the sender account.
		Receiver string  `json:"receiver"` // ID of the receiver account.
		Amount   float32 `json:"amount"`   // Amount to be transferred.
	}

	// Read and parse the JSON payload from the request.
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Get the user session from the request.
	session := app.GetSession(r)
	if session == nil {
		app.errorJSON(w, errAuthFailed, http.StatusUnauthorized)
		return
	}

	// Retrieve the sender account from the repository.
	senderAccount, err := app.Repo.GetAccountByID(requestPayload.Sender)
	if err != nil {
		if err == sql.ErrNoRows {
			app.errorJSON(w, errAccountNotFound)

		} else {
			app.errorJSON(w, err)
		}
		return
	}
	if senderAccount == nil {
		app.errorJSON(w, errAccountNotFound)
		return
	}

	// User tries to transfer money from an account which owns to another user
	if senderAccount.UserID != session.UserID {
		app.errorJSON(w, errAuthFailed, http.StatusUnauthorized)
		return
	}
	// Retrieve the receiver account from the repository.
	receiver, _ := app.Repo.GetAccountByID(requestPayload.Receiver)

	// Calculate the transfer fee based on whether the receiver is in the same bank.
	fee := func(sender, receiver *data.Account) float32 {
		if receiver == nil {
			sameBank = false
		}

		if !sameBank {
			return DifferentBankFee
		}
		return 0
	}(senderAccount, receiver)

	balance := senderAccount.Balance
	totalPayment := fee + requestPayload.Amount

	// Check if the sender has sufficient balance to cover the total payment.
	if totalPayment > balance {
		app.errorJSON(w, errInsufficentBalance)
		return
	}

	// Deduct the total payment from the sender's balance.
	senderAccount.Balance -= totalPayment

	// If the receiver is in the same bank, credit the amount to the receiver's balance.
	if sameBank {
		receiver.Balance += requestPayload.Amount
		// Update the receiver's account in the repository.
		err = app.Repo.UpdateAccount(receiver)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	}

	// Update the sender's account in the repository.
	err = app.Repo.UpdateAccount(senderAccount)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Prepare the response payload.
	payload := jsonResponse{
		Error:   false,
		Message: "Başarılı bir şekilde gönderildi",
	}

	// Record the transaction in the database.
	transaction := data.Transaction{
		SenderID:   requestPayload.Sender,
		ReceiverID: requestPayload.Receiver,
		Amount:     requestPayload.Amount,
		Type:       senderAccount.AccountType,
	}
	_, err = app.Repo.InsertTransaction(transaction)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Log the transaction asynchronously via RabbitMQ.
	go func() {
		d := map[string]interface{}{
			"sender_account":   requestPayload.Sender,
			"receiver_account": requestPayload.Receiver,
			"amount":           requestPayload.Amount,
			"fee":              fee,
			"account_type":     senderAccount.AccountType,
		}
		data, _ := json.Marshal(d)
		payload := LogPayload{
			Name: "log",
			Data: string(data),
		}
		_ = app.logEventViaRabbit(payload)
	}()

	// Send a confirmation email asynchronously via RabbitMQ.
	go func() {
		payload := MailPayload{
			From:    "transfer@teknasyon.com",
			To:      "test@teknasyon.com",
			Subject: "Para Gönderme İşlemi",
			Message: "Para gönderme işleminiz başarılı bir şekilde tamamlanmıştır.",
		}
		_ = app.mailEventViaRabbit(payload)
	}()

	// Send the response to the client.
	app.writeJSON(w, http.StatusAccepted, payload)
}

// CreateAccount handles the creation of a new account.
func (app *Config) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		AccountName    string  `json:"account_name"`    // Name of the new account.
		AccountType    string  `json:"account_type"`    // Type of the new account.
		InitialBalance float32 `json:"initial_balance"` // Type of the new account.
	}

	session := app.GetSession(r)
	if session == nil {
		app.errorJSON(w, errAuthFailed, http.StatusUnauthorized)
		return
	}

	// Start a new transaction
	tx, err := app.Repo.Begin()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer tx.Rollback() // Ensure rollback in case of failure

	// Read and parse the JSON payload from the request.
	err = app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	log.Printf("the init balance is %+v \n ", requestPayload.InitialBalance)
	// Create a new account instance with the provided details.
	account := data.Account{
		UserID:      session.UserID,
		AccountName: requestPayload.AccountName,
		AccountType: requestPayload.AccountType,
		Balance:     requestPayload.InitialBalance,
	}

	// Insert the new account into the repository and get the account ID.
	accountID, err := app.Repo.InsertAccount(account)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Prepare the response payload with the new account ID.
	payload := jsonResponse{
		Error:   false,
		Message: "",
		Data: struct {
			AccountID string `json:"account_id"` // ID of the newly created account.
		}{AccountID: accountID},
	}

	// Send the response to the client.
	app.writeJSON(w, http.StatusAccepted, payload)
}

// ListAccounts handles the listing of all accounts for the authenticated user.
func (app *Config) ListAccounts(w http.ResponseWriter, r *http.Request) {

	// Get the user session from the request.
	session := app.GetSession(r)
	if session == nil {
		app.errorJSON(w, errAuthFailed, http.StatusUnauthorized)
		return
	}

	// Retrieve all accounts associated with the user's ID from the repository.
	accounts, err := app.Repo.GetAllAccounts(session.UserID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Prepare the response payload with the retrieved accounts.
	payload := jsonResponse{
		Error:   false,
		Message: "",
		Data:    accounts,
	}

	// Send the response to the client.
	app.writeJSON(w, http.StatusAccepted, payload)
}

// GetAccount returns the account with specified id.
func (app *Config) GetAccount(w http.ResponseWriter, r *http.Request) {

	// Get the user session from the request.
	session := app.GetSession(r)
	if session == nil {
		app.errorJSON(w, errAuthFailed, http.StatusUnauthorized)
		return
	}

	accountID := chi.URLParam(r, "accountID")

	// Retrieve all accounts associated with the user's ID from the repository.
	account, err := app.Repo.GetAccountByID(accountID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Prepare the response payload with the retrieved accounts.
	payload := jsonResponse{
		Error:   false,
		Message: "",
		Data:    account,
	}

	// Send the response to the client.
	app.writeJSON(w, http.StatusAccepted, payload)
}

// logEventViaRabbit logs an event using the logger-service. It makes the call by pushing the data to RabbitMQ.
func (app *Config) logEventViaRabbit(l LogPayload) error {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		return err
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	return nil
}

// logEventViaRabbit logs an event using the logger-service. It makes the call by pushing the data to RabbitMQ.
func (app *Config) mailEventViaRabbit(m MailPayload) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = app.pushToQueue("mail", string(data))
	if err != nil {
		return err
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	return nil
}

// pushToQueue pushes a message into RabbitMQ
func (app *Config) pushToQueue(name string, msg string) error {
	emitter, err := event.NewEventEmitter(app.RabbitClient)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, err := json.MarshalIndent(&payload, "", "\t")
	if err != nil {
		return err
	}

	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}
