package middleware

/**
func (db *DB) CheckAccountPermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetClaimsFromContext(r)
		if !ok {
			utils.ErrorMessage(w, http.StatusUnauthorized, utils.INVALID_TOKEN)
			return
		}

		user, err := db.GetUserById(claims.User_Id)
		if err != nil {
			utils.ErrorMessage(w, http.StatusPreconditionFailed, err)
		}

		accountNumber_str := r.URL.Query().Get("number")
		accountNumber, err := utils.StringToUint64(accountNumber_str)
		if err != nil {
			utils.ErrorMessage(w, http.StatusBadRequest, err)
			return
		}

		account, err := db.GetAccountByAccountNumber(accountNumber)
		if err != nil {
			utils.ErrorMessage(w, http.StatusForbidden, err)
			return
		}
		if !user.HasAccount(account.ID) {
			utils.ErrorMessage(w, http.StatusNotFound, utils.ACCOUNT_NOT_FOUND)
			return
		}
		ctx := context.WithValue(r.Context(), "account", account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

*/
