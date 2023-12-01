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

	// defer cursor.Close(context.TODO())
	if err := cursor.All(context.TODO(), result); err != nil {
		log.Println("Ошибка при декодировании в FindReturnDecoded:", err)
		return err
	}

	return nil
}

// func main() {
// 	var myData []MyStruct

// 	filter := bson.M{"field1": "value"}
// 	collName := "collection_name"

// 	if err := Find(filter, collName, &myData); err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println(myData)
// }
