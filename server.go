package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
)


type Server struct {
	db *gorm.DB
}

var server = &Server{}

func reportHttpError(w http.ResponseWriter, err error) {
	http.Error(w, "", http.StatusInternalServerError)
	log.Println(err)
}

func writeJsonResponse(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("content-type", "application/json")
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		log.Println(fmt.Errorf("Error http writing response: %v", err))
	}
}

func SetupDB(driver, conn string) (*gorm.DB, error) {
	db, err := gorm.Open(driver, conn)
  if err != nil {
    return nil, err
  }

  db.AutoMigrate(&User{})
	db.AutoMigrate(&Transfer{})

	return db, nil
}

func InitTestData(w http.ResponseWriter, r *http.Request) {
	server.db.Delete(&User{})
	user1 := &User{FirstName: "Jean", LastName: "Carter", Email: "jcarter@gmail.com", Points: 500}
	user2 := &User{FirstName: "Bill", LastName: "Jones", Email: "bjones@gmail.com", Points: 600}
	server.db.Create(user1)
	server.db.Create(user2)

	server.db.Delete(&Transfer{})
	transfer1 := &Transfer{ UserID: user1.ID, Amount: 1000 }
	transfer2 := &Transfer{ UserID: user1.ID, Amount: -500 }
	transfer3 := &Transfer{ UserID: user2.ID, Amount: 600 }
	server.db.Create(transfer1)
	server.db.Create(transfer2)
	server.db.Create(transfer3)

	ListUsers(w, r)
}


func main() {
	db, err := SetupDB("sqlite3", "/tmp/button.db")
	if err != nil {
		panic(err)
	}
	server.db = db
	defer db.Close()

  r := mux.NewRouter()
	r.HandleFunc("/init-test-data", InitTestData)

	r.HandleFunc("/users/{id}/transfers", ListTransfers).Methods("GET")
	r.HandleFunc("/users/{id}/transfers", CreateTransfer).Methods("POST")
	r.HandleFunc("/users/{id}", GetUser).Methods("GET")
  r.HandleFunc("/users", ListUsers).Methods("GET")
	r.HandleFunc("/users", CreateUser).Methods("POST")

  log.Fatal(http.ListenAndServe(":8080", r))
}
