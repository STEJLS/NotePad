package main

import (
	"crypto/md5"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Error struct {
	Status int
	Text   string
}

type User struct {
	ID       bson.ObjectId  `bson:"_id,omitempty"`
	Login    string         `bson:"login"`
	Password [md5.Size]byte `bson:"password"`
}

type Note struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Author      string        `bson:"author"`
	Header      string        `bson:"header"`
	Text        string        `bson:"text"`
	CreatedDate time.Time     `bson:"createdDate"`
}
