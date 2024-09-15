package controllers

import "net/http"

func GetTransactions(w http.ResponseWriter, r *http.Request)            {}
func GetTransactionById(w http.ResponseWriter, r *http.Request)         {}
func GetTransactionsFromAccount(w http.ResponseWriter, r *http.Request) {}
func DepositToAccount(w http.ResponseWriter, r *http.Request)           {}
func WithdrawFromAccount(w http.ResponseWriter, r *http.Request)        {}
func TransferBetweenAccounts(w http.ResponseWriter, r *http.Request)    {}
