# Key-Value Memtable Server with MongoDB Backend written in Go

Inspired to write this after reading roughly 200 pages of Designing Data Intensive Applications and wanting to practice my Go.

Features:

    RESTful API: Exposes an endpoint (/send) to receive and store key-value pairs via HTTP POST requests.

    MongoDB Backend: Utilizes MongoDB as the backend storage for key-value pairs. The in-memory table is flushed to MongoDB when the size threshold is exceeded.

    Size-Based Flushing: The server flushes the in-memory table to MongoDB when the cumulative size of stored data reaches a predefined maximum.

Prerequisites:

    Go installed on your system.
    MongoDB server running and accessible.

Configuration:

    MongoDB Connection: Update the MongoDB connection URI in the main function.
    `opts := options.Client().ApplyURI("yoururl").SetServerAPIOptions(serverAPI)`

    Server IP Address: Set the IP address to bind the server in the handleRequests function.
    log.Fatal(http.ListenAndServe("yourip", myRouter))

    Maximum Size: Adjust the max field in the AppState struct to set the maximum cumulative size before flushing to MongoDB.
    state := &AppState{
    memTable: make(map[string][3]string),
    size:     0.0,
    max:      2,  // Set your desired maximum size
    }

Usage:
  Run the server:
      set IP and preffered Max size
      Ensure MongoDB is running with Mongo database ip set up
      go run main.go
  Access Service::
      Using Postman:
        POST IP:PORT/send?key=pair&database=collection

Testing:

  

Contributing

    Feel free to contribute to enhance or customize the functionality of this key-value storage server. Open issues and pull requests are welcome.
    Anyone is also free to utilize this in their project just please dm me because it would be cool to see someone else using my code.
