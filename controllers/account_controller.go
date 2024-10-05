package controllers

import (
	"github.com/mathis-k/bank-api/middleware"
	"github.com/mathis-k/bank-api/models"
	"github.com/mathis-k/bank-api/utils"
	"net/http"
)

func (s *APIServer) GetAccounts(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaimsFromContext(r)
	if !ok {
		utils.ErrorMessage(w, http.StatusUnauthorized, utils.INVALID_TOKEN)
		return
	}

	accounts, err := s.Database.GetAccountsFromUser(claims.User_Id)
	if err != nil {
		utils.ErrorMessage(w, http.StatusPreconditionFailed, err)
		return
	}
	utils.ResponseMessage(w, http.StatusOK, accounts)
}
func (s *APIServer) GetAccountByNumber(w http.ResponseWriter, r *http.Request) {
	account := r.Context().Value("account").(*models.Account)
	utils.ResponseMessage(w, http.StatusOK, account)
}
func (s *APIServer) CreateAccount(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaimsFromContext(r)
	if !ok {
		utils.ErrorMessage(w, http.StatusUnauthorized, utils.INVALID_TOKEN)
		return
	}

	account, err := s.Database.CreateAccount()
	if err != nil {
		utils.ErrorMessage(w, http.StatusInternalServerError, err)
		return
	}

	err = s.Database.AddAccountToUser(claims.User_Id, account.ID)
	if err != nil {
		utils.ErrorMessage(w, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseMessage(w, http.StatusCreated, account)

}
func (s *APIServer) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaimsFromContext(r)
	if !ok {
		utils.ErrorMessage(w, http.StatusUnauthorized, utils.INVALID_TOKEN)
		return
	}

	account := r.Context().Value("account").(*models.Account)

	err := s.Database.RemoveAccountFromUser(claims.User_Id, account.ID)
	if err != nil {
		utils.ErrorMessage(w, http.StatusInternalServerError, err)
		return
	}

	err = s.Database.DeleteAccount(account.ID)
	if err != nil {
		utils.ErrorMessage(w, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseMessage(w, http.StatusNoContent, `{"Success": "Account deleted"}`)
}
