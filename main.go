package main

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Name struct {
	Name string `json:"name"`
}

type Input struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type DatabaseRecord struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name,omitempty"`
	Text string             `bson:"text,omitempty"`
}

var client *mongo.Client
var collection *mongo.Collection
var err error

func main() {
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://192.168.1.2:27017"))
	if err != nil {
		panic(err)
	}

	collection = client.Database("testing").Collection("users")

	insert := http.HandlerFunc(handleInsert)
	http.Handle("/insert", insert)

	get := http.HandlerFunc(handleGet)
	http.Handle("/get", get)

	http.ListenAndServe(":8080", nil)
}

func handleInsert(w http.ResponseWriter, r *http.Request) {
	var text Input
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userText := bson.D{{"name", text.Name}, {"text", text.Text}}
	// insert the bson object using InsertOne()
	_, err = collection.InsertOne(context.TODO(), userText)
	// check for errors in the insertion
	if err != nil {
		panic(err)
	}

	resp := make(map[string]string)
	resp["Name"] = text.Name
	resp["Text"] = text.Text
	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
	return
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	var name Name
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// retrieving the first document that matches the filter
	var result DatabaseRecord
	// check for errors in the finding
	if err = collection.FindOne(context.TODO(), bson.M{"name": name.Name}).Decode(&result); err != nil {
		panic(err)
	}

	var text Input
	text.Name = result.Name
	text.Text = result.Text
	resp := make(map[string]string)
	resp["Name"] = text.Name
	resp["Text"] = text.Text
	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
	return
}
