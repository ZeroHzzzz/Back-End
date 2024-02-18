package models

type Cause struct {
	UserID string `bson:"UserID"`
	Msg    string `bson:"Msg"`
}
