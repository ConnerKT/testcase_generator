package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

// Define a struct to match the structure of MongoDB documents
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

	//Loading our .env file, and returning error if anything else happens
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	//Our URI for the database
	uri := os.Getenv("DB_URI")
	//Api key for the OPENAI API
	apiKey := os.Getenv("OPENAI_API_KEY")
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

	problem := Problem{

		Description: `<p>Given an array of integers <code>nums</code>&nbsp;and an integer <code>target</code>, return <em>indices of the two numbers such that they add up to <code>target</code></em>.</p>\n\n<p>You may assume that each input would have <strong><em>exactly</em> one solution</strong>, and you may not use the <em>same</em> element twice.</p>\n\n<p>You can return the answer in any order.</p>\n\n<p>&nbsp;</p>\n<p><strong class=\"example\">Example 1:</strong></p>\n\n<pre>\n<strong>Input:</strong> nums = [2,7,11,15], target = 9\n<strong>Output:</strong> [0,1]\n<strong>Explanation:</strong> Because nums[0] + nums[1] == 9, we return [0, 1].\n</pre>\n\n<p><strong class=\"example\">Example 2:</strong></p>\n\n<pre>\n<strong>Input:</strong> nums = [3,2,4], target = 6\n<strong>Output:</strong> [1,2]\n</pre>\n\n<p><strong class=\"example\">Example 3:</strong></p>\n\n<pre>\n<strong>Input:</strong> nums = [3,3], target = 6\n<strong>Output:</strong> [0,1]\n</pre>\n\n<p>&nbsp;</p>\n<p><strong>Constraints:</strong></p>\n\n<ul>\n\t<li><code>2 &lt;= nums.length &lt;= 10<sup>4</sup></code></li>\n\t<li><code>-10<sup>9</sup> &lt;= nums[i] &lt;= 10<sup>9</sup></code></li>\n\t<li><code>-10<sup>9</sup> &lt;= target &lt;= 10<sup>9</sup></code></li>\n\t<li><strong>Only one valid answer exists.</strong></li>\n</ul>\n\n<p>&nbsp;</p>\n<strong>Follow-up:&nbsp;</strong>Can you come up with an algorithm that is less than <code>O(n<sup>2</sup>)</code><font face=\"monospace\">&nbsp;</font>time complexity?`,

		FunctionSignature: "func twoSum(nums []int, target int) []int {\n\n\n}",
	}
// 	jsonData := []byte(`[
//     {
//         "nums": [2,7,11,15],
//         "target": 9
//     },
//     {
//         "nums": [ 3,2,4],
//         "target": 6
//     },
//     {
//         "nums": [
//             3,3
//         ],
//         "target": 6
//     }
// ]`)

	prompt := fmt.Sprintf(`Given the following leetcode problem:`, problem, `generate input values that match the functionSignatures`)
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
	f, err := os.Create("test.txt")
	l, err := f.WriteString(resp.Choices[0].Message.Content)
	fmt.Println(l)
	if err != nil {
		fmt.Println(err)
        f.Close()
		return
	}
	fmt.Println(resp.Choices[0].Message.Content)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

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
		// fmt.Println("Python:", result.FunctionSignatures.Python)
		// fmt.Println("JavaScript:", result.FunctionSignatures.Javascript)
		// fmt.Println("--------------------------------")
		// fmt.Println(result)
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Cursor error: %v", err)
	}
}
