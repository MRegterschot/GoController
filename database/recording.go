package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Recording struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
	Type string             `bson:"type"`
	Mode string             `bson:"mode"`
	Maps []MapRecords       `bson:"maps"`

	CreatedAt primitive.DateTime  `bson:"createdAt"`
	UpdatedAt primitive.DateTime  `bson:"updatedAt"`
	DeletedAt *primitive.DateTime `bson:"deletedAt,omitempty"`
}

type MapRecords struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	MapID       primitive.ObjectID `bson:"mapID"`
	MatchRounds []MatchRound       `bson:"matchRounds,omitempty"`
	Rounds      []Round            `bson:"rounds,omitempty"`
	Finishes    []PlayerFinish     `bson:"finishes,omitempty"`
}

type MatchRound struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	RoundNumber int                `bson:"roundNumber"`
	Teams       []Team             `bson:"teams"`
}

type Round struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	RoundNumber int                `bson:"roundNumber"`
	Players     []PlayerRound      `bson:"players"`
}

type Team struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	TeamID      int                `bson:"teamID"`
	Name        string             `bson:"name"`
	Points      int                `bson:"points"`
	TotalPoints int                `bson:"totalPoints"`
	Players     []PlayerRound      `bson:"players"`
}

type PlayerRound struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	PlayerID    *primitive.ObjectID `bson:"playerID"`
	Login       string              `bson:"login"`
	AccountId   string              `bson:"accountId"`
	Points      int                 `bson:"points"`
	TotalPoints int                 `bson:"totalPoints"`
	Time        int                 `bson:"time"`
	Checkpoints []int               `bson:"checkpoints"`
}

type PlayerFinish struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty"`
	PlayerID    *primitive.ObjectID `bson:"playerID"`
	Login       string              `bson:"login"`
	AccountId   string              `bson:"accountId"`
	Time        int                 `bson:"time"`
	Checkpoints []int               `bson:"checkpoints"`
	Timestamp   primitive.DateTime  `bson:"timestamp"`
}

func NewRecording(recording Recording) Recording {
	now := primitive.NewDateTimeFromTime(time.Now())

	newRecording := Recording{
		ID:        primitive.NewObjectID(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	CopyRecording(recording, &newRecording)

	return newRecording
}

func (m *Recording) Update(recording Recording) {
	now := primitive.NewDateTimeFromTime(time.Now())
	CopyRecording(recording, m)
	m.UpdatedAt = now
}

func (m *Recording) Delete() {
	now := primitive.NewDateTimeFromTime(time.Now())
	m.DeletedAt = &now
}

func CopyRecording(src Recording, dest *Recording) {
	dest.Name = src.Name
	dest.Mode = src.Mode
	dest.Type = src.Type
	dest.Maps = src.Maps
}

func GetRecordingByID(ctx context.Context, id primitive.ObjectID) (Recording, error) {
	var recording Recording
	err := GetCollection("recordings").FindOne(ctx, bson.M{"_id": id}).Decode(&recording)
	return recording, err
}

func InsertRecording(ctx context.Context, recording Recording) (*mongo.InsertOneResult, error) {
	return GetCollection("recordings").InsertOne(ctx, recording)
}

func UpdateRecording(ctx context.Context, recording Recording) (*mongo.UpdateResult, error) {
	return GetCollection("recordings").UpdateOne(ctx, bson.M{"_id": recording.ID}, bson.M{"$set": recording})
}
