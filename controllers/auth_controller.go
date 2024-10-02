package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mathis-k/bank-api/middleware"
	"github.com/mathis-k/bank-api/models"
	"github.com/mathis-k/bank-api/utils"
	"net/http"
)

func (s *APIServer) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userRequest models.UserRequest

	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		utils.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	if err := models.ValidateUserRequest(&userRequest); err != nil {
		utils.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}
	user, err := s.Database.CreateUser(&userRequest)
	if err != nil {
		if errors.Is(err, utils.EMAIL_ALREADY_EXISTS) {
			utils.ErrorMessage(w, http.StatusConflict, utils.EMAIL_ALREADY_EXISTS)
			return
		}
		utils.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	utils.ResponseMessage(w, http.StatusCreated, user)
}

func (s *APIServer) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userLogin models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		utils.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	if err := models.ValidateUserLogin(&userLogin); err != nil {
		utils.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	user, err := s.Database.LoginUser(&userLogin)
	if err != nil {
		if errors.Is(err, utils.INVALID_CREDENTIALS) {
			utils.ErrorMessage(w, http.StatusUnauthorized, utils.INVALID_CREDENTIALS)
			return
		}
		utils.ErrorMessage(w, http.StatusBadRequest, err)
		return
	}

	token, err := middleware.GenerateUserJWT(user.ID)
	if err != nil {
		utils.ErrorMessage(w, http.StatusInternalServerError, err)
		return
	}

	info := fmt.Sprintf("Welcome %s %s, your session/token is valid for %s. "+
		"Please make sure to enter your token in the authorization-header: "+
		"Authorization: Bearer <token>",
		user.FirstName, user.LastName, utils.FormatDuration(middleware.EXPIRATION_TIME_USER))
	msg := map[string]string{"Message": info, "token": token}
	utils.ResponseMessage(w, http.StatusOK, msg)
}
