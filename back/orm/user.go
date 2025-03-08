package orm

import (
	"crypto/rand"
	"encoding/binary"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

var (
	cryptoRand = rand.Reader
)

type User struct {
	ID         uint64    `json:"id" gorm:"primaryKey" faker:"-"`
	Name       string    `json:"name" gorm:"type:varchar(255);unique" faker:"name"`
	Email      string    `json:"-" gorm:"type:varchar(255);unique,omitempty" faker:"email"`
	CreatedAt  time.Time `json:"-" gorm:"type:timestamp" faker:"-"`
	RandomSeed uint64    `json:"-" faker:"-" gorm:"type:bigint"`

	Attempts []Attempt `json:"-" gorm:"foreignKey:UserID" faker:"-"`
}

func FakeUsers(n uint) []User {
	users := make([]User, n)
	faker.FakeData(&users, options.WithRandomMapAndSliceMinSize(n), options.WithRandomMapAndSliceMaxSize(n))

	for i := range users {
		buf := make([]byte, 8)
		if _, err := cryptoRand.Read(buf); err != nil {
			panic(err)
		}
		users[i].RandomSeed = binary.BigEndian.Uint64(buf) % 9223372036854775807
	}

	faker.ResetUnique()
	return users
}
