package main

import (
    "context"
    "fmt"
    "log"
    "time"
	"os"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
	"testcase_generator/models"

)

type FunctionSignature struct {
    Name     string `bson:"name"`
    Language string `bson:"language"`
    Value    string `bson:"value"`
}

type TestCase struct {
    ID     string                 `bson:"id"`
    Input  map[string]interface{} `bson:"input"`
    Output []interface{}          `bson:"output"`
}

func Seed() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	uri := os.Getenv("DB_URI")

    clientOptions := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB!")

    collection := client.Database("your_database").Collection("your_collection")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Fetch challenges
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        log.Fatal(err)
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var challenge models.Challenge
        err := cursor.Decode(&challenge)
        if err != nil {
            log.Fatal(err)
        }

        // Check if functionSignatures is in the correct format
        var updatedFunctionSignatures []FunctionSignature

        switch v := models.Challenge.FunctionSignatures.(type) {
        case map[string]interface{}:
            for lang, val := range v {
                if strVal, ok := val.(string); ok {
                    updatedFunctionSignatures = append(updatedFunctionSignatures, FunctionSignature{
                        Name:     extractFunctionName(strVal),
                        Language: lang,
                        Value:    strVal,
                    })
                }
            }
        case []interface{}:
            for _, sig := range v {
                if sigMap, ok := sig.(map[string]interface{}); ok {
                    name := sigMap["name"].(string)
                    language := sigMap["language"].(string)
                    value := sigMap["value"].(string)
                    updatedFunctionSignatures = append(updatedFunctionSignatures, FunctionSignature{
                        Name:     name,
                        Language: language,
                        Value:    value,
                    })
                }
            }
        }

        // If the updatedFunctionSignatures is different, update the document
        if len(updatedFunctionSignatures) > 0 && len(updatedFunctionSignatures) != len(challenge.FunctionSignatures.([]interface{})) {
            challenge.FunctionSignatures = updatedFunctionSignatures

            filter := bson.M{"_id": challenge.ID}
            update := bson.M{"$set": bson.M{"functionSignatures": updatedFunctionSignatures}}

            _, err := collection.UpdateOne(ctx, filter, update)
            if err != nil {
                log.Fatal(err)
            }
            fmt.Printf("Updated challenge with ID: %s\n", challenge.ID.Hex())
        }
    }

    if err := cursor.Err(); err != nil {
        log.Fatal(err)
    }

    // Close the connection once no longer needed
    if err := client.Disconnect(ctx); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connection to MongoDB closed.")
}

func extractFunctionName(functionValue string) string {
    // Dummy implementation to extract the function name from the function value
    // This function needs to be implemented according to your requirements
    if len(functionValue) == 0 {
        return ""
    }
    // Find the first occurrence of '(' and extract the function name
    idx := 0
    for i, char := range functionValue {
        if char == '(' {
            idx = i
            break
        }
    }
    funcSignature := functionValue[:idx]
    var funcName string
    fmt.Sscanf(funcSignature, "def %s", &funcName)
    return funcName
}
