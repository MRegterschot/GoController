package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MatchRecording struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
	Mode string             `bson:"mode"`
	Maps []MapRecords       `bson:"maps"`

	CreatedAt primitive.DateTime  `bson:"createdAt"`
	UpdatedAt primitive.DateTime  `bson:"updatedAt"`
	DeletedAt *primitive.DateTime `bson:"deletedAt,omitempty"`
}

type MapRecords struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	MapID  primitive.ObjectID `bson:"mapID"`
	Rounds []Round            `bson:"rounds"`
}

type Round struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	RoundNumber int                `bson:"roundNumber"`
	Teams       []Team             `bson:"teams"`
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
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	PlayerID    *primitive.ObjectID `bson:"playerID"`
	Login       string             `bson:"login"`
	Points      int                `bson:"points"`
	TotalPoints int                `bson:"totalPoints"`
	Time        int                `bson:"time"`
	Checkpoints []int              `bson:"checkpoints"`
}

func NewMatchRecording(matchRecording MatchRecording) MatchRecording {
	now := primitive.NewDateTimeFromTime(time.Now())

	newMatchRecording := MatchRecording{
		ID:        primitive.NewObjectID(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	CopyMatchRecording(matchRecording, &newMatchRecording)

	return newMatchRecording
}

func (m *MatchRecording) Update(matchRecording MatchRecording) {
	now := primitive.NewDateTimeFromTime(time.Now())
	CopyMatchRecording(matchRecording, m)
	m.UpdatedAt = now
}

func (m *MatchRecording) Delete() {
	now := primitive.NewDateTimeFromTime(time.Now())
	m.DeletedAt = &now
}

func CopyMatchRecording(src MatchRecording, dest *MatchRecording) {
	dest.Name = src.Name
	dest.Mode = src.Mode
	dest.Maps = src.Maps
}

func GetMatchRecordingByID(ctx context.Context, id primitive.ObjectID) (MatchRecording, error) {
	var matchRecording MatchRecording
	err := GetCollection("matchRecordings").FindOne(ctx, bson.M{"_id": id}).Decode(&matchRecording)
	return matchRecording, err
}

func InsertMatchRecording(ctx context.Context, matchRecording MatchRecording) (*mongo.InsertOneResult, error) {
	return GetCollection("matchRecordings").InsertOne(ctx, matchRecording)
}

func UpdateMatchRecording(ctx context.Context, matchRecording MatchRecording) (*mongo.UpdateResult, error) {
	return GetCollection("matchRecordings").UpdateOne(ctx, bson.M{"_id": matchRecording.ID}, bson.M{"$set": matchRecording})
}
