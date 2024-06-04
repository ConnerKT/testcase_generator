// This is the seed file to alter the databases functionSequence to the correct format
// This currently doesn't work, why? Don't ask me

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testcase_generator/models"
)

func functionSignature() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get MongoDB URI from environment variable
	uri := os.Getenv("DB_TESTING")
	//get chat api key
	apiKey := os.Getenv("OPENAI_API_KEY")

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

	// Access the db and collection
	collection := client.Database("testingLeet").Collection("problems")

	// Find challenges
	ctx := context.Background()
	// finds all the problems (idk how)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	//wtf is a cursor
	defer cursor.Close(ctx)
	f, err := os.OpenFile("practiceFS.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()


	for cursor.Next(context.TODO()) {
		var result models.Challenge
		if err := cursor.Decode(&result); err != nil {
			log.Fatalf("Error decoding document: %v", err)
		}

		fmt.Printf("Decoded document: %+v\n", result)

		prompt := fmt.Sprintf(`Given the following leetcode problem description: %s,`, result.Description, `generate a function signature for the problem:`, result.Title, `formatted like this, I want JavaScript and Python:`, `{
			{
				Name:    "title",
				Language: "python",
				Value: "the function signature",
			},
			{
				Name:    "title",
				Language: "javascript",
				Value: "the function signature",
			},
		}`)

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
			log.Fatalf("ChatCompletion error: %v", err)
		}
		if len(resp.Choices) == 0 {
			log.Fatalf("No choices returned from OpenAI")
		}

		content := resp.Choices[0].Message.Content
		if _, err := f.WriteString(content + "\n"); err != nil {
			log.Fatalf("Error writing to file: %v", err)
		}

		fmt.Println(content)
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Cursor error: %v", err)
	}
}

func main() {
	functionSignature()
}
