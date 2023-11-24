package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dataBase *mongo.Database

func Connect() {
	uri := "mongodb://" + os.Getenv("LOGIN") + ":" + os.Getenv("PASS") + "@" + os.Getenv("SERVER")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal("Ошибка подключения к базе данный =>", err)
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Fatal("Ping провален =>", err)
	}

	log.Println("База данных подключена упешно!")

	dataBase = client.Database(os.Getenv("BASE"))

	return
}

func InsertIfNotExists(filter, update primitive.M, collName string) *mongo.UpdateResult {
	opts := options.Update().SetUpsert(true)
	ctx := context.TODO()
	result, err := dataBase.Collection(collName).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Println("=InsertIfNotExists=", err)
	}
	return result
	// if result.MatchedCount != 0 {
	// 	fmt.Println("matched and replaced an existing document")
	// }
	// if result.UpsertedCount != 0 {
	// 	fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	// }
}

func UpdateIfExists(filter, update primitive.M, collName string) *mongo.UpdateResult {
	opts := options.Update().SetUpsert(false)
	result, err := dataBase.Collection(collName).UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Println("=UpdateIfExists=", err)
	}
	return result
}

func CountDocuments(filter primitive.M, collName string) int64 {
	ctx := context.TODO()
	itemCount, err := dataBase.Collection(collName).CountDocuments(ctx, filter)
	if err != nil {
		log.Println("=2671f1=", err)
	}
	return itemCount

}

func Find(filter primitive.M, collName string) *mongo.Cursor {
	cursor, err := dataBase.Collection(collName).Find(context.TODO(), filter)
	if err != nil {
		log.Println("=Find=", err)
	}
	return cursor
}

func FindSkip(filter primitive.M, collName string, skip, limit int) *mongo.Cursor {
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))
	cursor, err := dataBase.Collection(collName).Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Println("=Find=", err)
	}
	return cursor
}

func FindOne(filter primitive.M, collName string) *mongo.SingleResult {
	ctx := context.TODO()
	cursor := dataBase.Collection(collName).FindOne(ctx, filter)
	return cursor
}

// func Find(filter, sort bson.M, limit int64, collName string) (*mongo.Cursor, error) {
// 	findOptions := options.Find()
// 	findOptions.SetSort(sort)
// 	findOptions.SetLimit(limit)
// 	return DataBase.Collection(collName).Find(context.TODO(), filter, findOptions)
// }

// func FindOneAndUpdate(filter, update bson.M, upsert bool, collName string) *mongo.SingleResult {
// 	after := options.After
// 	opt := options.FindOneAndUpdateOptions{
// 		ReturnDocument: &after,
// 		Upsert:         &upsert,
// 	}
// 	return a := Collection(collName).FindOneAndUpdate(context.TODO(), filter, update, &opt)
// }
