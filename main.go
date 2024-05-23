package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type Challenge struct {
	ID                 primitive.ObjectID `bson:"_id"`
	Title              string             `bson:"title"`
	Difficulty         string             `bson:"difficulty"`
	Description        string             `bson:"description"`
	Link               string             `bson:"link"`
	FunctionSignatures struct {
		Python     string `bson:"python"`
		Javascript string `bson:"javascript"`
	} `bson:"functionSignatures"`
	TestCases []struct {
		ID    string `bson:"id"`
		Input struct {
			Nums   []int `bson:"nums"`
			Target int   `bson:"target"`
		} `bson:"input"`
		Output []int `bson:"output"`
	} `bson:"testCases"`
}
type Problem struct {
	Description       string
	FunctionSignature string
}

func main() {
	err := godotenv.Load()
	if err != nil {
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

	// Open the file in append mode before the loop
	f, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	for cursor.Next(context.TODO()) {
		var result Challenge
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatalf("Error decoding document: %v", err)
		}

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

	
		_, err = f.WriteString(resp.Choices[0].Message.Content + "\n")
		if err != nil {
			log.Fatalf("Error writing to file: %v", err)
		}

		fmt.Println(resp.Choices[0].Message.Content)
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Cursor error: %v", err)
	}
}
