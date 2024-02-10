package models

type Cause struct {
	UserId string `bson:"_id"`
	Msg    string `bson:"msg"`
}
