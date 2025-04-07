package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var mapsCollection = "maps"

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

	// Fields from nadeo API
	Submitter    string             `bson:"submitter"`
	Timestamp    primitive.DateTime `bson:"timestamp"`
	FileUrl      string             `bson:"fileUrl"`
	ThumbnailUrl string             `bson:"thumbnailUrl"`

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
	dest.Submitter = src.Submitter
	dest.Timestamp = src.Timestamp
	dest.FileUrl = src.FileUrl
	dest.ThumbnailUrl = src.ThumbnailUrl
}

func GetMapByID(ctx context.Context, id primitive.ObjectID) (Map, error) {
	var mapInfo Map
	err := GetCollection(mapsCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&mapInfo)
	return mapInfo, err
}

func GetMapByUId(ctx context.Context, uid string) (Map, error) {
	var mapInfo Map
	err := GetCollection(mapsCollection).FindOne(ctx, bson.M{"uid": uid}).Decode(&mapInfo)
	return mapInfo, err
}

func GetMapsByUIds(ctx context.Context, uids []string) ([]Map, error) {
	var maps []Map
	filter := bson.M{"uid": bson.M{"$in": uids}}

	cursor, err := GetCollection(mapsCollection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &maps); err != nil {
		return nil, err
	}

	return maps, nil
}

func InsertMap(ctx context.Context, mapInfo Map) (*mongo.InsertOneResult, error) {
	return GetCollection(mapsCollection).InsertOne(ctx, mapInfo)
}

func InsertMaps(ctx context.Context, maps []Map) (*mongo.InsertManyResult, error) {
	docs := make([]any, len(maps))
	for i, mapInfo := range maps {
		docs[i] = mapInfo
	}
	return GetCollection(mapsCollection).InsertMany(ctx, docs)
}

func UpdateMap(ctx context.Context, mapInfo Map) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": mapInfo.ID}
	update := bson.M{
		"$set": mapInfo,
	}

	return GetCollection(mapsCollection).UpdateOne(ctx, filter, update)
}
