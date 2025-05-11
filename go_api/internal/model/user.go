package model

import (
	"time"
)

type User struct {
	ID             string    `bson:"_id"`
	Username       string    `bson:"username"`
	Email          string    `bson:"email"`
	FirstName      string    `bson:"firstName"`
	LastName       string    `bson:"lastName"`
	PasswordHashed string    `bson:"passwordHashed"`
	CreatedDate    time.Time `bson:"createdDate"`
	UpdatedDate    time.Time `bson:"UpdatedDate"`
}
