package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Challenge struct {
    ID                 primitive.ObjectID `bson:"_id"`
    Title              string             `bson:"title"`
    Difficulty         string             `bson:"difficulty"`
    Description        string             `bson:"description"`
    Link               string             `bson:"link"`
    FunctionSignatures []FunctionSignature `bson:"functionSignatures"`
}

type FunctionSignature struct {
    Name     string `bson:"name"`
    Language string `bson:"language"`
    Value    string `bson:"value"`
}
