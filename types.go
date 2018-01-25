package main

import (
	"crypto/md5"

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
