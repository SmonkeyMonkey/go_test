package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/testproject/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users/{page}", usersHandler).Methods(http.MethodGet)
	r.HandleFunc("/users", createHandler).Methods(http.MethodPost)
	r.HandleFunc("/update/{id}", updateHandler).Methods(http.MethodPut)

	http.ListenAndServe(":8080", r)
}

func usersHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	currentPage, err := strconv.Atoi(vars["page"])
	if err != nil {
		log.Fatal("error conv to int")
	}
	users := models.Users(currentPage)
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(rw, users)
}
func createHandler(rw http.ResponseWriter, r *http.Request) {
	var user models.User
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &user)
	result,err :=user.Create(user)
	if err != nil{
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}else {
		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(result)
		rw.WriteHeader(http.StatusOK)
		rw.Write(response)
	}
}
func updateHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	objId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
	u := models.User{}
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &u)
	result, err := u.Update(objId, u)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	} else {
		rw.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(result)
		rw.WriteHeader(http.StatusOK)
		rw.Write(response)
	}

}
