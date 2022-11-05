package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

type Direction struct {
	Name string `json:"direction"`
}

var NameArray [5]string

func main() {
	NameArray[0] = "Damian"
	NameArray[1] = "Ann"
	NameArray[2] = "Elizabeth"
	NameArray[3] = "Jacob"
	NameArray[4] = "Mark"
	handler := http.HandlerFunc(handleRequest)
	http.Handle("/directions", handler)
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var direction Direction
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&direction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := make(map[string]string)
	resp["Name"] = NameArray[rand.Intn(5)]
	resp["Direction"] = direction.Name
	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
	return
}
