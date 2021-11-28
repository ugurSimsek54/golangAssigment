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
	Startdate time.Time             `json:"startdate,omitempty" bson:"startdate,omitempty"`
	Enddate  time.Time             `json:"enddate,omitempty" bson:"enddate,omitempty"`
    Mincount int				`json:"mincount,omitempty" bson:"mincount,omitempty"`
	Maxcount int				`json:"maxcount,omitempty" bson:"maxcount,omitempty"`
}

type RecordRequest struct {
    ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Key        string `json:"key,omitempty" bson:"key,omitempty"`
	Createdat time.Time             `json:"createdat,omitempty" bson:"createdat,omitempty"`
    Totalcount int				`json:"totalcount,omitempty" bson:"totalcount,omitempty"`
}

var client *mongo.Client

func CreateRecordEndpoint(response http.ResponseWriter, request *http.Request){
	response.Header().Add("content-type","application/json")
	var record RecordRequest
	_=json.NewDecoder(request.Body).Decode(&record)
	collection := client.Database("getircase-study").Collection("records")
	ctx,_:=context.WithTimeout(context.Background(),10*time.Second)
	result,_:=collection.InsertOne(ctx,record)
	json.NewEncoder(response).Encode(result)
}


func GetRecordsEndpoint(response http.ResponseWriter, request *http.Request) {
	
	response.Header().Set("content-type", "application/json")
	var recordRequest []RecordRequest
	var record Records
	_=json.NewDecoder(request.Body).Decode(&record)
	fmt.Println(record.Mincount)
	fmt.Println(record.Maxcount)
	collection := client.Database("getircase-study").Collection("records")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"createdat":bson.M{"$gte":primitive.NewDateTimeFromTime(record.Startdate),"$lte":primitive.NewDateTimeFromTime(record.Enddate)},"totalcount":bson.M{"$gt":record.Mincount,"$lt":record.Maxcount}})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var recordreq RecordRequest
		cursor.Decode(&recordreq)
		recordRequest = append(recordRequest, recordreq)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(recordRequest)
}
func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getircase-study?retryWrites=true")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/record", CreateRecordEndpoint).Methods("POST")
	router.HandleFunc("/records", GetRecordsEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}
