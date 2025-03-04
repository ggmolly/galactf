package orm

import (
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

type User struct {
	ID        uint64    `json:"id" gorm:"primaryKey" faker:"-"`
	Name      string    `json:"name" gorm:"type:varchar(255);unique" faker:"name"`
	Email     string    `json:"email" gorm:"type:varchar(255);unique" faker:"email"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp" faker:"-"`

	Attempts []Attempt `json:"attempts" gorm:"foreignKey:UserID" faker:"-"`
}

func FakeUsers(n uint) []User {
	users := make([]User, n)
	faker.FakeData(&users, options.WithRandomMapAndSliceMinSize(n), options.WithRandomMapAndSliceMaxSize(n))
	faker.ResetUnique()
	return users
}
