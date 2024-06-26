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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	uri := os.Getenv("DB_URI")
	apiKey := os.Getenv("OPENAI_API_KEY")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
		fmt.Println("Disconnected from MongoDB")
	}()
	fmt.Println("Connected to MongoDB!")

	collection := client.Database("leetcode").Collection("challenges")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("Error finding documents: %v", err)
	}
	defer cursor.Close(context.TODO())

	f, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

		prompt := fmt.Sprintf(`Given the following leetcode problem: %s, generate input values that match the functionSignatures %v`, result.Description, result.FunctionSignatures)
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
