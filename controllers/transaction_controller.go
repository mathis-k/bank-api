package controllers

import (
	"encoding/json"
	"github.com/mathis-k/bank-api/middleware"
	"github.com/mathis-k/bank-api/models"
	"github.com/mathis-k/bank-api/utils"
	"net/http"
)

func (s *APIServer) GetTransactions(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaimsFromContext(r)
	if !ok {
		utils.ErrorMessage(w, http.StatusUnauthorized, utils.INVALID_TOKEN)
		return
	}

	user, err := s.Database.GetUserById(claims.User_Id)
	if err != nil {
		utils.ErrorMessage(w, http.StatusPreconditionFailed, err)
		return
	}

	var transactions []models.Transaction
	for _, aId := range user.Accounts {
		transactions_, err := s.Database.GetTransactionsFromAccount(aId)
		if err != nil {
			utils.ErrorMessage(w, http.StatusInternalServerError, err)
			return
		}
		for _, t := range transactions_ {
			transactions = append(transactions, *t)
		}
	}

	utils.ResponseMessage(w, http.StatusOK, transactions)
}
func (s *APIServer) GetTransactionById(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaimsFromContext(r)
	if !ok {
		utils.ErrorMessage(w, http.StatusUnauthorized, utils.INVALID_TOKEN)
		return
	}

	transactionId := r.URL.Query().Get("id")

	transactions, err := s.Database.GetTransactionsFromUser(claims.User_Id)
	if err != nil {
		utils.ErrorMessage(w, http.StatusInternalServerError, err)
		return
	}

	for _, t := range transactions {
		if t.ID.Hex() == transactionId {
			utils.ResponseMessage(w, http.StatusOK, t)
			return
		}
	}

	utils.ErrorMessage(w, http.StatusNotFound, utils.TRANSACTION_NOT_FOUND)
}
func (s *APIServer) GetTransactionsFromAccount(w http.ResponseWriter, r *http.Request) {
	account := r.Context().Value("account").(models.Account)
	transactions, err := s.Database.GetTransactionsFromAccount(account.ID)
	if err != nil {
		utils.ErrorMessage(w, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseMessage(w, http.StatusOK, transactions)
}
func (s *APIServer) DepositToAccount(w http.ResponseWriter, r *http.Request) {
	account := r.Context().Value("account").(models.Account)

	var transactionRequest models.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&transactionRequest); err != nil {
		utils.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}
	transactionRequest.Type = "DEPOSIT"
	transactionRequest.ToAccountID = account.ID

	if err := models.ValidateTransactionRequest(&transactionRequest); err != nil {
		utils.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	transaction, err := s.Database.CreateTransaction(&transactionRequest)

	if err != nil {
		utils.ErrorMessage(w, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseMessage(w, http.StatusCreated, transaction)
}
func (s *APIServer) WithdrawFromAccount(w http.ResponseWriter, r *http.Request) {
	account := r.Context().Value("account").(models.Account)

	var transactionRequest models.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&transactionRequest); err != nil {
		utils.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	transactionRequest.Type = "WITHDRAW"
	transactionRequest.FromAccount = account.ID

	if err := models.ValidateTransactionRequest(&transactionRequest); err != nil {
		utils.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	transaction, err := s.Database.CreateTransaction(&transactionRequest)

	if err != nil {
		utils.ErrorMessage(w, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseMessage(w, http.StatusCreated, transaction)
}
func (s *APIServer) TransferBetweenAccounts(w http.ResponseWriter, r *http.Request) {}
