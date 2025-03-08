package orm

import (
	"errors"
	"time"

	"github.com/ggmolly/galactf/utils"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Challenge struct {
	ID          uint64         `json:"id" gorm:"primaryKey" faker:"-"`
	Name        string         `json:"name" gorm:"type:varchar(255);unique" faker:"word,unique,lang=eng,len=10"`
	Description string         `json:"description" gorm:"type:text" faker:"sentence"`
	Difficulty  uint8          `json:"difficulty" gorm:"type:smallint;index" faker:"number,boundary_start=0,boundary_end=6"`
	Categories  pq.StringArray `json:"categories" gorm:"type:varchar(24)[]" faker:"ctfCategory"`
	CreatedAt   time.Time      `json:"-" gorm:"autoCreateTime" faker:"-"`

	Attachments []Attachment `json:"attachments" gorm:"foreignKey:ChallengeID"`
	Attempts    []Attempt    `json:"attempts" gorm:"foreignKey:ChallengeID"`
}

type ChallengeStats struct {
	Challenge

	SolveRate float64 `json:"solve_rate" faker:"-" gorm:"-"`
}

func GetChallengeStats() ([]ChallengeStats, error) {
	var challenges []ChallengeStats

	// Load challenges along with their attachments
	err := GormDB.
		Table("challenges").
		Preload("Attachments").
		Find(&challenges).
		Error
	if err != nil {
		return nil, err
	}

	// Compute solve rate for each challenge
	for i := range challenges {
		var totalAttempts, successfulAttempts int64

		err := GormDB.Table("attempts").
			Where("challenge_id = ?", challenges[i].ID).
			Count(&totalAttempts).Error
		if err != nil {
			return nil, err
		}

		err = GormDB.Table("attempts").
			Where("challenge_id = ? AND success = true", challenges[i].ID).
			Count(&successfulAttempts).Error
		if err != nil {
			return nil, err
		}

		if totalAttempts > 0 {
			challenges[i].SolveRate = float64(successfulAttempts) / float64(totalAttempts)
		}
	}

	return challenges, nil
}

func GetChallengeStatsById(id int) (*ChallengeStats, error) {
	var result ChallengeStats
	err := GormDB.Table("challenges").
		Select("challenges.id, challenges.name, challenges.difficulty, challenges.categories, challenges.description, "+
			"COUNT(attempts.id) FILTER(WHERE attempts.success = true) * 1.0 / NULLIF(COUNT(attempts.id), 0) AS solve_rate").
		Joins("LEFT JOIN attempts ON attempts.challenge_id = challenges.id").
		Where("challenges.id = ?", id).
		Group("challenges.id").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetChallengeById(id int) (*Challenge, error) {
	var result Challenge
	if err := GormDB.First(&result, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.RestStatusFactory(nil, fiber.StatusNotFound, "challenge not found")
		} else {
			return nil, utils.RestStatusFactory(nil, fiber.StatusInternalServerError, "error fetching challenge")
		}
	}
	return &result, nil
}

func GetFakeChallenges(n uint) []Challenge {
	chals := make([]Challenge, n)
	faker.FakeData(&chals, options.WithRandomMapAndSliceMinSize(n), options.WithRandomMapAndSliceMaxSize(n))
	faker.ResetUnique()
	return chals
}
