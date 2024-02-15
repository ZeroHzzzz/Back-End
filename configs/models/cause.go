package models

type Cause struct {
	UserID string `bson:"_id"`
	Msg    string `bson:"Msg"`
}
