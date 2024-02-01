package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"unsafe"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AppState holds the application state
type AppState struct {
	memTable map[string][3]string
	size     float32
	max      int
	coll     mongo.Database
	mutex    sync.Mutex
}

func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return string(localAddress.IP)
}

func storePair(w http.ResponseWriter, r *http.Request, state *AppState, client *mongo.Client) {
	fmt.Print("Hello, ")
	var temp float32
	tempSlice := make([]string, 0)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Iterate over all form key-value pairs
	var counter int = 0
	for key, values := range r.Form {
		tempSlice = append(tempSlice, key)
		for _, value := range values {
			tempSlice = append(tempSlice, value)
			if counter < 2 {
				temp = temp + float32(uint(unsafe.Sizeof(value)))
			}
			key = key
		}
		counter = counter + 1

	}

	if state.size+temp > float32(state.max) {
		flushTable(state.memTable, tempSlice[2], tempSlice[3], client) // magic numbers for database and collection to insert into
		state.memTable[tempSlice[0]] = [3]string{tempSlice[1], tempSlice[2], tempSlice[3]}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Data stored successfully"))
		fmt.Print("Hello, ")
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Data stored not!"))
		state.memTable[tempSlice[0]] = [3]string{tempSlice[1], tempSlice[2], tempSlice[3]}
		fmt.Print("bye, ")
	}

}

func flushTable(mem_table map[string][3]string, data string, collect string, client *mongo.Client) {
	for key, values := range mem_table {
		document := bson.M{
			"key":   key,
			"value": bson.A{values[0], values[1], values[2]}, 
		}

		collection := client.Database(data).Collection(collect)
		_, err := collection.InsertOne(context.TODO(), document)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func handleRequests(state *AppState, client *mongo.Client) {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		storePair(w, r, state, client)
	}).Methods("POST")

	log.Fatal(http.ListenAndServe("yourip", myRouter))
}

func main() {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("yoururl").SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	//getting database

	state := &AppState{
		memTable: make(map[string][3]string),
		size:     0.0,
		max:      2,
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	handleRequests(state, client)
}
