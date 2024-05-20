package main

import (
    "context"
    "fmt"
    "log"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Define a struct to match the structure of MongoDB documents
type Challenge struct {
    FunctionSignatures struct {
        Python     string `bson:"python"`
        Javascript string `bson:"javascript"`
    } `bson:"functionSignatures"`
}

func main() {

	//Our URI for the database
    uri := "mongodb+srv://codecrew_User:1234@cluster0.nnbdzbq.mongodb.net/leetcode?retryWrites=true&w=majority"

	// We provide two variables, one for the client and one for errors
	// We connect and provide a default context, and apply the URI
	//Destructuring our connect 

    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
    if err != nil {
        log.Fatalf("Error connecting to MongoDB: %v", err)
    }
	//A function used when the main function ends (we disconnect)
    defer func() {
        if err := client.Disconnect(context.TODO()); err != nil {
            log.Fatalf("Error disconnecting from MongoDB: %v", err)
        }
        fmt.Println("Disconnected from MongoDB")
    }()
	//connection line
    fmt.Println("Connected to MongoDB!")
	// setting our client database, if it doesn't exist, create it
    collection := client.Database("leetcode").Collection("challenges")

    // we perform the query here, finding all the documents
	//we create a pointer
    cursor, err := collection.Find(context.TODO(), bson.D{})
    if err != nil {
        log.Fatalf("Error finding documents: %v", err)
    }
	//if the cursor is done (loop) we close it
    defer cursor.Close(context.TODO())

	// we loop checking each document, and creating a copy using the challenge struct
    for cursor.Next(context.TODO()) {
		var result Challenge

		// Decode the next document from the cursor into the 'result' variable.
		// If an error occurs during decoding, assign it to the 'err' variable.
		err := cursor.Decode(&result)
		
		// Check if there was an error during decoding.
		if err != nil {
			// If an error occurred, log a fatal error message with the error details and exit the program.
			log.Fatalf("Error decoding document: %v", err)
		}

        // Access function signatures for Python and JavaScript
        fmt.Println("Python:", result.FunctionSignatures.Python)
        fmt.Println("JavaScript:", result.FunctionSignatures.Javascript)
        fmt.Println("--------------------------------")
    }

    if err := cursor.Err(); err != nil {
        log.Fatalf("Cursor error: %v", err)
    }
}
