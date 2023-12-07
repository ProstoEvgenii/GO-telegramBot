package server

import (
	"GO-chatModeratorTg/pages"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var router = map[string]func(http.ResponseWriter, *http.Request){
	"UserAuth":       pages.AuthHandler,
	"Admins":         pages.AdminsHandler,
	"Mentions":       pages.MentionsHandler,
	"Urls":           pages.UrlsHandler,
	"ForbiddenWords": pages.ForbiddenWordsHandler,
}

func HandleRequest(rw http.ResponseWriter, request *http.Request) {

	path := strings.Split(request.URL.Path, "/api/")

	handler, exists := router[path[1]]

	if exists {
		handler(rw, request)
	} else {
		log.Println("Не найден handler => ", path[1])
		// Обработка случая, когда маршрут не найден
		http.NotFound(rw, request)
	}
}

func PostToApi(route string, requestBody []byte) (io.ReadCloser, error) {
	base := "https://api.telegram.org/bot" + os.Getenv("token") + "/" + route
	res, err := http.Post(base, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return nil, err
	}
	return res.Body, nil
}

func GetToApi(route string) (io.ReadCloser, error) {
	base := "https://api.telegram.org/bot" + os.Getenv("token") + "/" + route
	res, err := http.Get(base)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return nil, err
	}
	return res.Body, nil
}
func Start(host string) {
	log.Println("=ba58dc=")
	http.HandleFunc("/", HandleRequest)
	http.ListenAndServe(host, nil)
}
