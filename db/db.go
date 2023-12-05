package db

import (
	"context"
	"fmt"
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

func InsertIfNotExists(filter, update primitive.M, collName string, upsert bool) bool {
	opts := options.Update().SetUpsert(upsert)
	ctx := context.TODO()
	result, err := dataBase.Collection(collName).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Println("=InsertIfNotExists=", err)
	}

	if result.UpsertedCount != 0 {
		return true
	}
	return false
}

func DeleteDocument(filter primitive.M, collName string) bool {
	ctx := context.TODO()
	res, err := dataBase.Collection(collName).DeleteOne(ctx, filter)
	if err != nil {
		log.Println("=305e43=", err)
	}
	if res.DeletedCount != 0 {
		return true
	}
	return false
}

func InsertOne(document interface{}, collName string) {
	ctx := context.TODO()
	result, err := dataBase.Collection(collName).InsertOne(ctx, document)
	if err != nil {
		log.Println("=Find=", err)
	}
	fmt.Printf("Document inserted: %v\n", result)
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

func FindReturnCursor(filter primitive.M, collName string, result interface{}) *mongo.Cursor {
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

func FindReturnDecoded(filter interface{}, collName string, result interface{}) error {
	cursor, err := dataBase.Collection(collName).Find(context.TODO(), filter)
	if err != nil {
		log.Println("Ошибка cursor в FindReturnDecoded:", err)
		return err
	}

	if err := cursor.All(context.TODO(), result); err != nil {
		log.Println("Ошибка при декодировании в FindReturnDecoded:", err)
		return err
	}
	return nil
}
