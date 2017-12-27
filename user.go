package main

import(
  "encoding/json"
  "net/http"
  "strconv"

  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
)

type User struct {
  ID 				uint 	  `json:"id"`
  Email 		string 	`json:"email"`
  FirstName string 	`json:"first_name"`
  LastName 	string 	`json:"last_name"`
  Points 		int 		`json:"points"`
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
  var users []User
  server.db.Find(&users)
  writeJsonResponse(w, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
  user := FindUserFromRequest(server.db, w, r)
  if user == nil {
    return
  }
  writeJsonResponse(w, *user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
  var user User
  err := json.NewDecoder(r.Body).Decode(&user)
  if err != nil {
  reportHttpError(w, err)
  return
  }
  if err = server.db.Create(&user).Error; err != nil {
    reportHttpError(w, err)
    return
  }
  writeJsonResponse(w, user)
}

func FindUserFromRequest(db *gorm.DB, w http.ResponseWriter, r *http.Request) *User {
  params := mux.Vars(r)
  id, err := strconv.Atoi(params["id"])

  if err != nil {
    reportHttpError(w, err)
    return nil
  }

  var user User
  if db.First(&user, id).RecordNotFound() {
    http.Error(w, "user not found", http.StatusNotFound)
    return nil
  }

  if db.Error != nil {
    reportHttpError(w, err)
    return nil
  }

  return &user
}
