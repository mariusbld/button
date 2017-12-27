package main

import(
  "encoding/json"
  "net/http"
)

type Transfer struct {
  ID 				uint 	  `json:"id"`
  UserID    uint    `json:"user_id"`
  Amount 		int 		`json:"amount"`
}

func ListTransfers(w http.ResponseWriter, r *http.Request) {
  user := FindUserFromRequest(server.db, w, r)
  if user == nil {
    return
  }

  var transfers []Transfer
  if err := server.db.Where(&Transfer{ UserID: user.ID }).Find(&transfers).Error; err != nil {
    reportHttpError(w, err)
    return
  }

  writeJsonResponse(w, transfers)
}

func CreateTransfer(w http.ResponseWriter, r *http.Request) {
  tx := server.db.Begin()
  defer tx.Rollback()

  user := FindUserFromRequest(tx, w, r)
  if user == nil {
    return
  }

  var transfer Transfer
  err := json.NewDecoder(r.Body).Decode(&transfer)
  if err != nil {
  reportHttpError(w, err)
  return
  }

  if user.ID != transfer.UserID {
    http.Error(w, "user id mismatch", http.StatusBadRequest)
    return
  }

  if user.Points + transfer.Amount < 0 {
    http.Error(w, "not enough points", http.StatusPaymentRequired)
    return
  }

  user.Points += transfer.Amount

  tx.Create(&transfer)
  tx.Save(&user)

  tx.Commit()
  writeJsonResponse(w, transfer)
}
