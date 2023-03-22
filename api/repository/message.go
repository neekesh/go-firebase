package repository

import (
	"context"
	"fmt"
	"letschat/infrastructure"
	"letschat/models"
	"letschat/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

func NewMessageRepository(db infrastructure.Database, logger infrastructure.Logger) MessageRepository {
	return MessageRepository{
		db:     db,
		logger: logger,
	}
}

func (c MessageRepository) Create(message models.MessageM) error {
	usersCollection := c.db.DB.Collection("messages")
	_, err := usersCollection.InsertOne(context.TODO(), message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (c MessageRepository) Delete(id string) error {
	usersCollection := c.db.DB.Collection("messages")
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	_, err := usersCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (c MessageRepository) FindAll(pagination utils.Pagination, roomId string) (*[]models.MessageM, string, error) {

	var messages []models.MessageM

	usersCollection := c.db.DB.Collection("messages")

	filter := bson.M{"room_id": roomId}

	opt := options.Find()

	if pagination.Cursor != "" {
		objID, err := primitive.ObjectIDFromHex(pagination.Cursor)

		if err != nil {
			fmt.Println(err)
			return nil, "", err
		}
		filter = bson.M{"room_id": roomId, "_id": bson.M{"$gt": objID}}
	}
	opt.SetLimit(int64(pagination.Limit))

	cur, err := usersCollection.Find(context.TODO(), filter, opt)

	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	fmt.Println("Cursor value:", cur)
	for cur.Next(context.TODO()) {
		var message models.MessageM
		err := cur.Decode(&message)
		if err != nil {
			fmt.Println(err)
			return nil, "", err
		}
		messages = append(messages, message)
	}

	if err := cur.Err(); err != nil {
		fmt.Println(err)
		return nil, "", err
	}
	// Get the next cursor string by getting the _id of the last document
	var nextCursor string

	if len(messages) > 0 {
		nextCursor = messages[len(messages)-1].MessageId.Hex()

	}

	return &messages, nextCursor, nil
}
