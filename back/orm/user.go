package orm

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/ggmolly/galactf/cache"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

var (
	cryptoRand = rand.Reader
)

type User struct {
	ID         uint64    `json:"id" gorm:"primaryKey" faker:"-"`
	Name       string    `json:"name" gorm:"type:varchar(255);unique" faker:"name"`
	Email      string    `json:"-" gorm:"type:varchar(255);unique,omitempty,index" faker:"email"`
	CreatedAt  time.Time `json:"-" gorm:"type:timestamp" faker:"-"`
	RandomSeed uint64    `json:"-" faker:"-" gorm:"type:bigint"`

	Attempts []Attempt `json:"-" gorm:"foreignKey:UserID" faker:"-"`
}

type RedisUser struct {
	ID         uint64 `json:"id" gorm:"primaryKey" faker:"-"`
	Name       string `json:"name" gorm:"type:varchar(255);unique" faker:"name"`
	RandomSeed uint64 `json:"random_seed" faker:"-" gorm:"type:bigint"`
}

func GenerateRandomSeed() (uint64, error) {
	buf := make([]byte, 8)
	if _, err := cryptoRand.Read(buf); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(buf) % 9223372036854775807, nil
}

func FakeUsers(n uint) []User {
	var err error
	users := make([]User, n)
	faker.FakeData(&users, options.WithRandomMapAndSliceMinSize(n), options.WithRandomMapAndSliceMaxSize(n))

	for i := range users {
		if users[i].RandomSeed, err = GenerateRandomSeed(); err != nil {
			panic(err)
		}
	}

	faker.ResetUnique()
	return users
}

func InvalidateChallengesCache(userID uint64) {
	cache.InvalidateKey(fmt.Sprintf(cache.ChallengesCacheKey, userID))
}
