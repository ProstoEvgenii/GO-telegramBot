package pages

import (
	"GO-chatModeratorTg/db"
	"GO-chatModeratorTg/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson"
)

func MentionsHandler(rw http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		totalFound := db.CountDocuments(bson.M{"type": "mention"}, "whiteList")
		params := new(models.Params)
		filter := bson.M{
			"type": "mention",
		}
		limitPerPage := 15
		page := 1

		if err := schema.NewDecoder().Decode(params, request.URL.Query()); err != nil {
			log.Println("=Params schema Error Database=", err)
		}

		_, exists := AuthUsers[params.UUID]
		if !exists {
			return
		}

		if params.Page != 0 {
			page = params.Page
		}
		skip := limitPerPage * (page - 1)
		cursor := db.FindSkip(filter, "whiteList", skip, limitPerPage)
		var usersSlice []models.WhiteList
		if err := cursor.All(context.TODO(), &usersSlice); err != nil {
			log.Println("Cursor All Error Database", err)
			rw.Write([]byte("{}"))
			return
		}
		response := models.Admins_response{
			Records:    usersSlice,
			TotalFound: totalFound,
			Page:       page,
		}
		dataBaseJson, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error:", err)
			rw.Write([]byte("{}"))
			return
		}
		rw.Write(dataBaseJson)
		return
	}
	if request.Method == "POST" {

		rw.Write([]byte("Привет"))
	}
}
