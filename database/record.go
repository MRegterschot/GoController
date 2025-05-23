package database

import (
	"context"
	"time"

	"github.com/MRegterschot/GoController/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var recordsCollection = "records"

type Record struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Login  string             `bson:"login"`
	Time   int                `bson:"time"`
	MapUId string             `bson:"mapUid"`

	CreatedAt primitive.DateTime  `bson:"createdAt"`
	UpdatedAt primitive.DateTime  `bson:"updatedAt"`
	DeletedAt *primitive.DateTime `bson:"deletedAt,omitempty"`
}

func NewRecord(record Record) Record {
	now := primitive.NewDateTimeFromTime(time.Now())

	newRecord := Record{
		ID:        primitive.NewObjectID(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	CopyRecord(record, &newRecord)

	return newRecord
}

func CopyRecord(src Record, dest *Record) {
	dest.Login = src.Login
	dest.Time = src.Time
	dest.MapUId = src.MapUId
}

func GetRecordByID(ctx context.Context, id primitive.ObjectID) (Record, error) {
	var record Record
	err := GetCollection(recordsCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&record)
	return record, err
}

func GetRecordByLogin(ctx context.Context, login string) (Record, error) {
	var record Record
	err := GetCollection(recordsCollection).FindOne(ctx, bson.M{"login": login}).Decode(&record)
	return record, err
}

func GetRecordsByLogin(ctx context.Context, login string) ([]Record, error) {
	var records []Record
	cursor, err := GetCollection(recordsCollection).Find(ctx, bson.M{"login": login})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var record Record
		if err := cursor.Decode(&record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func GetRecordsByMapUId(ctx context.Context, mapUId string) ([]Record, error) {
	var records []Record
	cursor, err := GetCollection(recordsCollection).Find(ctx, bson.M{"mapUid": mapUId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var record Record
		if err := cursor.Decode(&record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func InsertRecord(ctx context.Context, record Record) (*mongo.InsertOneResult, error) {
	return GetCollection(recordsCollection).InsertOne(ctx, record)
}

func (r *Record) ToModel(dest *models.Record) {
	dest.ID = r.ID.Hex()
	dest.Login = r.Login
	dest.Time = r.Time
	dest.MapUId = r.MapUId
	dest.CreatedAt = r.CreatedAt.Time()
	dest.UpdatedAt = r.UpdatedAt.Time()
	if r.DeletedAt != nil {
		t := r.DeletedAt.Time()
		dest.DeletedAt = &t
	}
}
