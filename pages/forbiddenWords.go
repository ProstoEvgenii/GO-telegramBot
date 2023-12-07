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

func ForbiddenWordsHandler(rw http.ResponseWriter, request *http.Request) {

	if request.Method == "GET" {
		totalFound := db.CountDocuments(bson.M{}, "forbiddenWords")
		params := new(models.Params)
		filter := bson.M{}
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
		if params.Seach != "" {
			filter = bson.M{
				"word": bson.M{"$regex": params.Seach, "$options": "i"},
			}
		}
		skip := limitPerPage * (page - 1)
		cursor := db.FindSkip(filter, "forbiddenWords", skip, limitPerPage)
		var forbiddenWords []models.ForbiddenWords
		if err := cursor.All(context.TODO(), &forbiddenWords); err != nil {
			log.Println("Cursor All Error Database", err)
			rw.Write([]byte("{}"))
			return
		}
		response := models.ForbiddenWords_response{
			Records:    forbiddenWords,
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
