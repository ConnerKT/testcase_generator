package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testcase_generator/models"
)

type FunctionSignature struct {
	Name     string `json:"name"`
	Language string `json:"language"`
	Value    string `json:"value"`
}

func functionSignature() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get MongoDB URI from environment variable
	uri := os.Getenv("DB_TESTING")
	apiKey := os.Getenv("OPENAI_API_KEY")

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting MongoDB: %v", err)
		}
	}()

	// Ping MongoDB
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// Access the db and collection
	collection := client.Database("testingLeet").Collection("problems")

	// Find challenges
	ctx := context.Background()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Fatalf("Error closing cursor: %v", err)
		}
	}()

	for cursor.Next(context.TODO()) {
		var result models.Challenge
		if err := cursor.Decode(&result); err != nil {
			log.Fatalf("Error decoding document: %v", err)
		}

		fmt.Printf("Decoded document: %+v\n", result)

		prompt := fmt.Sprintf(`Given the following leetcode problem description: %s,
generate a function signature for the problem: %s formatted as JSON like this, I want JavaScript and Python:
[
	{
		"name": "functionNameInPython",
		"language": "python",
		"value": "def functionNameInPython(params):"
	},
	{
		"name": "functionNameInJavaScript",
		"language": "javascript",
		"value": "function functionNameInJavaScript(params) { }"
	}
]`, result.Description, result.Title)

		var functionSignatures []FunctionSignature

		for attempts := 1; attempts <= 3; attempts++ {
			openClient := openai.NewClient(apiKey)
			resp, err := openClient.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model: openai.GPT3Dot5Turbo,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: prompt,
						},
					},
				},
			)
			if err != nil {
				log.Printf("ChatCompletion error: %v", err)
				time.Sleep(2 * time.Second)
				continue
			}
			if len(resp.Choices) == 0 {
				log.Printf("No choices returned from OpenAI, attempt %d", attempts)
				time.Sleep(2 * time.Second)
				continue
			}

			content := resp.Choices[0].Message.Content
			fmt.Printf("OpenAI response: %s\n", content)

			// Ensure the response is valid JSON
			content = strings.TrimSpace(content)
			if !strings.HasPrefix(content, "[") || !strings.HasSuffix(content, "]") {
				log.Printf("Invalid JSON format returned, attempt %d", attempts)
				time.Sleep(2 * time.Second)
				continue
			}

			err = json.Unmarshal([]byte(content), &functionSignatures)
			if err != nil {
				log.Printf("Error unmarshalling function signatures: %v, attempt %d", err, attempts)
				time.Sleep(2 * time.Second)
				continue
			}
			break
		}

		if functionSignatures == nil {
			log.Fatalf("Failed to get valid function signatures after 3 attempts")
		}

		// Update the document with the new field functionSignature
		filter := bson.M{"_id": result.ID}
		update := bson.M{
			"$set": bson.M{
				"functionSignature": functionSignatures,
			},
		}

		updateResult, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatalf("Error updating document: %v", err)
		}
		fmt.Printf("Updated document with _id: %v, Matched %v documents and updated %v documents.\n", result.ID, updateResult.MatchedCount, updateResult.ModifiedCount)
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Cursor error: %v", err)
	}
}

func main() {
	functionSignature()
}
