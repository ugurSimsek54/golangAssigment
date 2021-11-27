package main

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"encoding/json"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"github.com/gorilla/mux"
)

type Records struct {
    ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Key     string 				`json:"key,omitempty" bson:"key,omitempty"`
	Value 	string             	`json:"value,omitempty" bson:"value,omitempty"`
}

var client *mongo.Client

func CreateRecordEndpoint(response http.ResponseWriter, request *http.Request){
	response.Header().Add("content-type","application/json")
	var record Records
	_=json.NewDecoder(request.Body).Decode(&record)
	collection := client.Database("getircase-study-ınmemory").Collection("records")
	ctx,_:=context.WithTimeout(context.Background(),10*time.Second)
	result,_:=collection.InsertOne(ctx,record)
	json.NewEncoder(response).Encode(result)
}


func GetRecordsEndpoint(response http.ResponseWriter, request *http.Request) {
	
	response.Header().Set("content-type", "application/json")
	var records []Records
	var record Records
	_=json.NewDecoder(request.Body).Decode(&record)
	collection := client.Database("getircase-study-ınmemory").Collection("records")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"key":request.URL.Query()["key"][0]})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var recordreq Records
		cursor.Decode(&recordreq)
		records = append(records, recordreq)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(records)
}
func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/in-memory", CreateRecordEndpoint).Methods("POST")
	router.HandleFunc("/in-memory", GetRecordsEndpoint).Methods("GET")
	http.ListenAndServe(":12346", router)
}
