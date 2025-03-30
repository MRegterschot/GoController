package database

import (
	"context"
	"time"

	"github.com/MRegterschot/GoController/models"
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

func GetRecordings(ctx context.Context) ([]Recording, error) {
	cursor, err := GetCollection("recordings").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var recordings []Recording
	if err = cursor.All(ctx, &recordings); err != nil {
		return nil, err
	}

	return recordings, nil
}

func InsertRecording(ctx context.Context, recording Recording) (*mongo.InsertOneResult, error) {
	return GetCollection("recordings").InsertOne(ctx, recording)
}

func UpdateRecording(ctx context.Context, recording Recording) (*mongo.UpdateResult, error) {
	return GetCollection("recordings").UpdateOne(ctx, bson.M{"_id": recording.ID}, bson.M{"$set": recording})
}

func (r *Recording) ToModel(dest *models.Recording) {
	dest.ID = r.ID.Hex()
	dest.Name = r.Name
	dest.Type = r.Type
	dest.Mode = r.Mode
	dest.CreatedAt = r.CreatedAt.Time()
	dest.UpdatedAt = r.UpdatedAt.Time()
	if r.DeletedAt != nil {
		t := r.DeletedAt.Time()
		dest.DeletedAt = &t
	}

	dest.Maps = make([]models.MapRecords, 0, len(r.Maps))
	for i, mapRecord := range r.Maps {
		mapRecord.ToModel(&dest.Maps[i])
	}
}

func (m *MapRecords) ToModel(dest *models.MapRecords) {
	dest.ID = m.ID.Hex()

	dest.MatchRounds = make([]models.MatchRound, 0, len(m.MatchRounds))
	for i, matchRound := range m.MatchRounds {
		matchRound.ToModel(&dest.MatchRounds[i])
	}

	dest.Rounds = make([]models.Round, 0, len(m.Rounds))
	for i, round := range m.Rounds {
		round.ToModel(&dest.Rounds[i])
	}

	dest.Finishes = make([]models.PlayerFinish, 0, len(m.Finishes))
	for i, finish := range m.Finishes {
		finish.ToModel(&dest.Finishes[i])
	}
}

func (m *MatchRound) ToModel(dest *models.MatchRound) {
	dest.ID = m.ID.Hex()
	dest.RoundNumber = m.RoundNumber

	dest.Teams = make([]models.Team, 0, len(m.Teams))
	for i, team := range m.Teams {
		team.ToModel(&dest.Teams[i])
	}
}

func (m *Round) ToModel(dest *models.Round) {
	dest.ID = m.ID.Hex()
	dest.RoundNumber = m.RoundNumber

	dest.Players = make([]models.PlayerRound, 0, len(m.Players))
	for i, player := range m.Players {
		player.ToModel(&dest.Players[i])
	}
}

func (m *Team) ToModel(dest *models.Team) {
	dest.ID = m.ID.Hex()
	dest.TeamID = m.TeamID
	dest.Name = m.Name
	dest.Points = m.Points
	dest.TotalPoints = m.TotalPoints

	dest.Players = make([]models.PlayerRound, 0, len(m.Players))
	for i, player := range m.Players {
		player.ToModel(&dest.Players[i])
	}
}

func (m *PlayerRound) ToModel(dest *models.PlayerRound) {
	dest.ID = m.ID.Hex()
	dest.Login = m.Login
	dest.AccountId = m.AccountId
	dest.Points = m.Points
	dest.TotalPoints = m.TotalPoints
	dest.Time = m.Time
	dest.Checkpoints = m.Checkpoints
}

func (m *PlayerFinish) ToModel(dest *models.PlayerFinish) {
	dest.ID = m.ID.Hex()
	dest.Login = m.Login
	dest.AccountId = m.AccountId
	dest.Time = m.Time
	dest.Checkpoints = m.Checkpoints
	dest.Timestamp = m.Timestamp.Time()
}