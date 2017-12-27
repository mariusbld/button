marius/button

## Install:
```sh
go get github.com/mattn/go-sqlite3
go get github.com/gorilla/mux
go get -u github.com/jinzhu/gorm

git clone https://github.com/mariusbld/button.git $GOPATH/src/button
cd $GOPATH/src/button
go install
```
## Run:
```sh
button
```
The HTTP server will start with the default port 8080, using a default local SQLite db file at /tmp/button.db

## JSON Format: 
User:
```
{
  "id":20,
  "email":"jcarter@gmail.com",
  "first_name":"Jean"
  "last_name":"Carter",
  "points":100
}
```
Transfer:
```
{
  "id":19,
  "user_id":20,
  "amount":1000
}
```

## Rest API
List all users:
```http
GET /users
```
Get single user by id:
>- If the {id} is not found will return **404**
```http
GET /users/{id}
```
Create user:
>- If another user with the same {id} is present will return **500**
```http
POST /users
```
List all transfers for a specified user:
>- If the user is not found will return **404**
```http
GET /users/{id}/transfers
```
Create a new transfer
>- If the user is not found will return **404**
>- If there are not enough points, will return **402**
```http
POST /users/{id}/transfers
```

## Test API
Populate DB with sample test users and transfers:
```http
GET /init-test-data
```
