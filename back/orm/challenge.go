package orm

import (
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/lib/pq"
)

type Challenge struct {
	ID         uint64         `json:"id" gorm:"primaryKey" faker:"-"`
	Name       string         `json:"name" gorm:"type:varchar(255);unique" faker:"word,unique,lang=eng,len=10"`
	Difficulty uint8          `json:"difficulty" gorm:"type:smallint;index" faker:"number,boundary_start=0,boundary_end=6"`
	Categories pq.StringArray `json:"categories" gorm:"type:varchar(24)[]" faker:"ctfCategory"`
	CreatedAt  time.Time      `json:"-" gorm:"autoCreateTime" faker:"-"`
}

func GetChallenges() ([]Challenge, error) {
	var challenges []Challenge
	if err := GormDB.
		Order("difficulty ASC").
		Find(&challenges).
		Error; err != nil {
		return nil, err
	}
	return challenges, nil
}

func GetFakeChallenges(n uint) []Challenge {
	chals := make([]Challenge, n)
	faker.FakeData(&chals, options.WithRandomMapAndSliceMinSize(n), options.WithRandomMapAndSliceMaxSize(n))
	faker.ResetUnique()
	return chals
}
