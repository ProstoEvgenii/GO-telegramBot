package pages

import (
	"GO-chatModeratorTg/db"
	"GO-chatModeratorTg/models"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var AuthUsers = map[string]int64{}

func AuthHandler(rw http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		uploadAuth(rw, request)
		return
	}

	return
}

func uploadAuth(rw http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var params models.Auth
	err := decoder.Decode(&params)
	if err != nil {
		log.Println("Error parse post News => ", err)
		fmt.Fprintf(rw, "{\"error\":\"Неверные данные\"}")
		return
	}
	data := []byte(params.Password)
	hash := sha256.Sum256(data)

	hashString := hex.EncodeToString(hash[:])
	// log.Println("=c7f972=", hashString)
	var tmp interface{}
	filter := bson.M{
		"login":    params.Login,
		"password": hashString,
	}

	check := CheckInDB(tmp, filter, "auth")
	if !check {
		fmt.Fprintf(rw, "{\"error\":\" Не верные данные\"}")
		return
	} else {
		now := time.Now()
		timestamp := now.Unix()
		AuthUsers[params.UUID] = timestamp
		go CheckAuthUsers()
		fmt.Fprintf(rw, "{\"result\":\"Авторизация успешна\"}")
		return
	}

}

func CheckAuthUsers() {
	now := time.Now()
	timestampNow := now.Unix()
	for item := range AuthUsers {
		if timestampNow > AuthUsers[item]+(60*60*24) {
			delete(AuthUsers, item)
		}
	}
}
func CheckInDB(tmp interface{}, filter bson.M, collName string) bool {
	find := db.FindOne(filter, collName)
	if err := find.Decode(&tmp); err != nil {
		return false
	}
	return true
}
