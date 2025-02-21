package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Map struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Name           string             `bson:"name"`
	UId            string             `bson:"uid"`
	FileName       string             `bson:"fileName"`
	Author         string             `bson:"author"`
	AuthorNickname string             `bson:"authorNickname"`
	AuthorTime     int                `bson:"authorTime"`
	GoldTime       int                `bson:"goldTime"`
	SilverTime     int                `bson:"silverTime"`
	BronzeTime     int                `bson:"bronzeTime"`

	CreatedAt primitive.DateTime  `bson:"createdAt"`
	UpdatedAt primitive.DateTime  `bson:"updatedAt"`
	DeletedAt *primitive.DateTime `bson:"deletedAt,omitempty"`
}

func NewMap(m Map) Map {
	now := primitive.NewDateTimeFromTime(time.Now())

	newMap := Map{
		ID:        primitive.NewObjectID(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	CopyMap(m, &newMap)

	return newMap
}

func (m *Map) Update(mapInfo Map) {
	now := primitive.NewDateTimeFromTime(time.Now())
	CopyMap(mapInfo, m)
	m.UpdatedAt = now
}

func (m *Map) Delete() {
	now := primitive.NewDateTimeFromTime(time.Now())
	m.DeletedAt = &now
}

func CopyMap(src Map, dest *Map) {
	dest.Name = src.Name
	dest.UId = src.UId
	dest.FileName = src.FileName
	dest.Author = src.Author
	dest.AuthorNickname = src.AuthorNickname
	dest.AuthorTime = src.AuthorTime
	dest.GoldTime = src.GoldTime
	dest.SilverTime = src.SilverTime
	dest.BronzeTime = src.BronzeTime
}

func GetMapByID(ctx context.Context, id primitive.ObjectID) (Map, error) {
	var mapInfo Map
	err := GetCollection("maps").FindOne(ctx, bson.M{"_id": id}).Decode(&mapInfo)
	return mapInfo, err
}

func GetMapByUId(ctx context.Context, uid string) (Map, error) {
	var mapInfo Map
	err := GetCollection("maps").FindOne(ctx, bson.M{"uid": uid}).Decode(&mapInfo)
	return mapInfo, err
}

func InsertMap(ctx context.Context, mapInfo Map) (*mongo.InsertOneResult, error) {
	return GetCollection("maps").InsertOne(ctx, mapInfo)
}
