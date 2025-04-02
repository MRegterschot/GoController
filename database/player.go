package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var playersCollection = "players"

type Player struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Login    string             `bson:"login"`
	NickName string             `bson:"nickName"`
	Path     string             `bson:"path"`

	CreatedAt primitive.DateTime  `bson:"createdAt"`
	UpdatedAt primitive.DateTime  `bson:"updatedAt"`
	DeletedAt *primitive.DateTime `bson:"deletedAt,omitempty"`
}

func NewPlayer(player Player) Player {
	now := primitive.NewDateTimeFromTime(time.Now())

	newPlayer := Player{
		ID:        primitive.NewObjectID(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	CopyPlayer(player, &newPlayer)

	return newPlayer
}

func (p *Player) Update(player Player) {
	now := primitive.NewDateTimeFromTime(time.Now())
	CopyPlayer(player, p)
	p.UpdatedAt = now
}

func (p *Player) Delete() {
	now := primitive.NewDateTimeFromTime(time.Now())
	p.DeletedAt = &now
}

func CopyPlayer(src Player, dest *Player) {
	dest.Login = src.Login
	dest.NickName = src.NickName
	dest.Path = src.Path
}

func GetPlayerByID(ctx context.Context, id primitive.ObjectID) (Player, error) {
	var player Player
	err := GetCollection(playersCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&player)
	return player, err
}

func GetPlayerByLogin(ctx context.Context, login string) (Player, error) {
	var player Player
	err := GetCollection(playersCollection).FindOne(ctx, bson.M{"login": login}).Decode(&player)
	return player, err
}

func InsertPlayer(ctx context.Context, player Player) (*mongo.InsertOneResult, error) {
	return GetCollection(playersCollection).InsertOne(ctx, player)
}

func UpdatePlayer(ctx context.Context, player Player) (*mongo.UpdateResult, error) {
	return GetCollection(playersCollection).UpdateOne(ctx, bson.M{"_id": player.ID}, bson.M{"$set": player})
}
