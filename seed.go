package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/joho/godotenv"
    "testcase_generator/models"
)



func Seed() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Get MongoDB URI from environment variable
    uri := os.Getenv("DB_URI")

    // Connect to MongoDB
    clientOptions := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(context.TODO())

    // Ping MongoDB
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB!")

    // Access the collection
    collection := client.Database("leetcode").Collection("challenges")

    // Find challenges
    ctx := context.Background()
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        log.Fatal(err)
    }
    defer cursor.Close(ctx)

 // Iterate over challenges
for cursor.Next(ctx) {
    var challenge models.Challenge
    err := cursor.Decode(&challenge)
    if err != nil {
        log.Fatalf("Error decoding challenge document: %v", err)
    }
    
    // // Print out the challenge document for inspection
    // fmt.Printf("Challenge document: %+v\n", challenge)

    // Check if FunctionSignatures is an object
    if len(challenge.FunctionSignatures) == 1 {
        // Convert FunctionSignatures object to array
        updatedFunctionSignatures := []models.FunctionSignature{
            {
                Name:     challenge.FunctionSignatures[0].Name,
                Language: challenge.FunctionSignatures[0].Language,
                Value:    challenge.FunctionSignatures[0].Value,
            },
        }
        challenge.FunctionSignatures = updatedFunctionSignatures
    }

    // Update the document in the collection
    filter := bson.M{"_id": challenge.ID}
    update := bson.M{"$set": bson.M{"functionSignatures": challenge.FunctionSignatures}}
    _, err = collection.UpdateOne(ctx, filter, update)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Updated challenge with ID: %s\n", challenge.ID.Hex())
}



    // Check for errors during cursor iteration
    if err := cursor.Err(); err != nil {
        log.Fatal(err)
    }

    fmt.Println("Operation completed successfully.")
}

func main() {
	Seed();
}