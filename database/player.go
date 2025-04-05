package database

import (
	"context"
	"time"

	"github.com/MRegterschot/GoController/models"
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
	Roles    []string           `bson:"roles"`

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
	dest.Roles = src.Roles
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

func (p *Player) ToModel(dest *models.Player) {
	dest.ID = p.ID.Hex()
	dest.Login = p.Login
	dest.NickName = p.NickName
	dest.Path = p.Path
	dest.Roles = p.Roles
	dest.CreatedAt = p.CreatedAt.Time()
	dest.UpdatedAt = p.UpdatedAt.Time()
	if p.DeletedAt != nil {
		t := p.DeletedAt.Time()
		dest.DeletedAt = &t
	}
}
