package orm

import (
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/lib/pq"
)

type Challenge struct {
	ID          uint64         `json:"id" gorm:"primaryKey" faker:"-"`
	Name        string         `json:"name" gorm:"type:varchar(255);unique" faker:"word,unique,lang=eng,len=10"`
	Description string         `json:"description" gorm:"type:text" faker:"sentence"`
	Difficulty  uint8          `json:"difficulty" gorm:"type:smallint;index" faker:"number,boundary_start=0,boundary_end=6"`
	Categories  pq.StringArray `json:"categories" gorm:"type:varchar(24)[]" faker:"ctfCategory"`
	CreatedAt   time.Time      `json:"-" gorm:"autoCreateTime" faker:"-"`
}

type ChallengeSolveRate struct {
	Challenge

	SolveRate float64 `json:"solve_rate" faker:"-" gorm:"column:solve_rate"`
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

func GetChallengeSolveRate() ([]ChallengeSolveRate, error) {
	var result []ChallengeSolveRate

	err := GormDB.Table("challenges").
		Select("challenges.id, challenges.name, challenges.difficulty, challenges.categories, challenges.description, " +
			"COUNT(attempts.id) FILTER(WHERE attempts.success = true) * 1.0 / NULLIF(COUNT(attempts.id), 0) AS solve_rate").
		Joins("LEFT JOIN attempts ON attempts.challenge_id = challenges.id").
		Group("challenges.id").
		Order("challenges.difficulty ASC").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetFakeChallenges(n uint) []Challenge {
	chals := make([]Challenge, n)
	faker.FakeData(&chals, options.WithRandomMapAndSliceMinSize(n), options.WithRandomMapAndSliceMaxSize(n))
	faker.ResetUnique()
	return chals
}
