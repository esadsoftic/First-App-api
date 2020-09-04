package main

import (
	"encoding/json"
	"fmt"

	//"fmt"
	"github.com/globalsign/mgo/bson"
	"github.com/rs/cors"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"bitbucket/firstapp/config"

	"github.com/globalsign/mgo"
)

var configuration = config.Config{}
var dataService = PersonsDAO{}

type PersonsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "persons"
)

type Person struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	Name  string        `bson:"name" json:"name"`
	Lname string        `bson:"lname" json:"lname"`
	Age   int           `bson:"age" json:"age"`
}

type PersonsResponse struct {
	Persons []Person `bson:"person" json:"person"`
}

func (m *PersonsDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Println(m)
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

func (m *PersonsDAO) FindAll() ([]Person, error) {
	var persons []Person
	err := db.C(COLLECTION).Find(bson.M{}).All(&persons)
	return persons, err
}

func (m *PersonsDAO) FindById(id string) ([]Person, error) {
	var persons []Person
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&persons)
	return persons, err
}

func (m *PersonsDAO) Insert(persons *Person) error {

	err := db.C(COLLECTION).Insert(&persons)
	return err
}

func (m *PersonsDAO) Delete(id string) error {
	err := db.C(COLLECTION).RemoveId(bson.ObjectIdHex(id))
	return err
}

func (m *PersonsDAO) Update(persons Person) error {
	err := db.C(COLLECTION).UpdateId(persons.ID, &persons)
	return err
}

// AllPersonEndpoint is the function to all the exported endpoint
func AllPersonEndpoint(w http.ResponseWriter, r *http.Request) {
	persons, err := dataService.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, persons)
}

//p := Person{Name: "Josh", Age: 30}
//pr, err := json.Marshal(p)
//if err!= nil {
//	panic(err)
//}
//fmt.Println(string(pr))
//fmt.Fprintln(w, string(pr))
//}

// FindPersonEndpoint is the function to find the exported endpoint
func FindPersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	persons, err := dataService.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Person ID")
		return
	}
	respondWithJson(w, http.StatusOK, persons)
}

// CreatePersonEndpoint is the function to create the exported endpoint
func CreatePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	var person *Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request in CreatPersonEndpoint")
		return
	}
	person.ID = bson.NewObjectId()
	if err := dataService.Insert(person); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, person)

}

// UpdatePersonEndpoint is the function to update the exported endpoint
func UpdatePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var persons Person
	if err := json.NewDecoder(r.Body).Decode(&persons); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request UpdatePersonEndpoint")
		return
	}
	if err := dataService.Update(persons); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DeletePersonEndpoint is the function to delete the exported endpoint
func DeletePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Println(id)
	if err := dataService.Delete(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func init() {
	configuration.Read()

	dataService.Server = configuration.Server
	dataService.Database = configuration.Database
	dataService.Connect()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/people", AllPersonEndpoint).Methods("GET")
	r.HandleFunc("/people", CreatePersonEndpoint).Methods("POST")
	r.HandleFunc("/people/{id}", UpdatePersonEndpoint).Methods("PUT")
	r.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")
	r.HandleFunc("/people/{id}", FindPersonEndpoint).Methods("GET")

	c := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: false,
	}

	log.Println("We Good")
	if err := http.ListenAndServe(":8080", cors.New(c).Handler(r)); err != nil {
		log.Fatalf("Error from main: %v", err)
	}
}
