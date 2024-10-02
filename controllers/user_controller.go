package controllers

import (
	"encoding/json"
	"github.com/mathis-k/bank-api/middleware"
	"github.com/mathis-k/bank-api/models"
	"github.com/mathis-k/bank-api/utils"
	"net/http"
)

func (s *APIServer) GetUser(w http.ResponseWriter, r *http.Request) {
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

	utils.ResponseMessage(w, http.StatusOK, user)
}
func (s *APIServer) UpdateUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaimsFromContext(r)
	if !ok {
		utils.ErrorMessage(w, http.StatusUnauthorized, utils.INVALID_TOKEN)
		return
	}

	var userUpdate models.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&userUpdate); err != nil {
		utils.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	if err := models.ValidateUserUpdate(&userUpdate); err != nil {
		utils.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	user, err := s.Database.UpdateUser(claims.User_Id, &userUpdate)
	if err != nil {
		utils.ErrorMessage(w, http.StatusInternalServerError, err)
		return
	}

	utils.ResponseMessage(w, http.StatusOK, user)
}
