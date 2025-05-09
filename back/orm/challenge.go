package orm

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ggmolly/galactf/types"
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
	RevealAt    time.Time      `json:"-" faker:"-"`

	Attachments []Attachment `json:"attachments" gorm:"foreignKey:ChallengeID"`
	Attempts    []Attempt    `json:"-" gorm:"foreignKey:ChallengeID"`
}

type ChallengeStats struct {
	Challenge

	SolveRate     float64 `json:"solve_rate" faker:"-" gorm:"-"`
	Solved        bool    `json:"solved" faker:"-" gorm:"-"`
	Solvers       uint64  `json:"solvers" faker:"-" gorm:"-"`
	TotalAttempts int     `json:"total_attempts" faker:"-" gorm:"-"`
	RevealIn      uint64  `json:"reveal_in,omitempty" faker:"-" gorm:"-"` // relative time until reveal
}

func GetChallengeStats(userID uint64) ([]ChallengeStats, error) {
	var challenges []ChallengeStats

	// Load challenges along with their attachments
	err := GormDB.
		Table("challenges").
		Preload("Attempts").
		Preload("Attachments").
        Order("reveal_at ASC").
		Find(&challenges).
		Error
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	for i := range challenges {
		challenge := &challenges[i]

		solvedAttempts := 0
		for _, attempt := range challenge.Attempts {
			if attempt.Success {
				solvedAttempts++
			}
		}
		challenge.Solvers = uint64(solvedAttempts)
		challenge.TotalAttempts = len(challenge.Attempts)
		if challenge.TotalAttempts > 0 {
			challenge.SolveRate = float64(solvedAttempts) / float64(challenge.TotalAttempts)
		} else {
			challenge.SolveRate = 0
		}

		for _, attempt := range challenge.Attempts {
			if attempt.UserID == userID && attempt.Success {
				challenge.Solved = true
				break
			}
		}

		// Only serialize the reveal time if it's in the future
		if now.Before(challenge.RevealAt) {
			challenge.RevealIn = uint64(challenge.RevealAt.Sub(now).Seconds())
			// Censor other informations
			challenge.Name = ""
			challenge.Description = ""
			challenge.Categories = []string{}
			challenge.Attachments = []Attachment{}
			challenge.Attempts = []Attempt{}
		}
	}

	return challenges, nil
}

func GetChallengeStatsById(id int, userID uint64) (*ChallengeStats, error) {
	var challenge ChallengeStats

	err := GormDB.
		Table("challenges").
		Preload("Attempts").
		Preload("Attachments").
		Where("challenges.id = ?", id).
		First(&challenge).
		Error
	if err != nil {
		return nil, err
	}

	challenge.TotalAttempts = len(challenge.Attempts)
	solvedAttempts := 0
	for _, attempt := range challenge.Attempts {
		if attempt.Success {
			solvedAttempts++
		}
	}

	if challenge.TotalAttempts > 0 {
		challenge.SolveRate = float64(solvedAttempts) / float64(challenge.TotalAttempts)
	} else {
		challenge.SolveRate = 0
	}
	challenge.Solvers = uint64(solvedAttempts)
	for _, attempt := range challenge.Attempts {
		if attempt.UserID == userID && attempt.Success {
			challenge.Solved = true
			break
		}
	}

	return &challenge, nil
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

func SendFirstBlood(chal *Challenge, solver *User) {
    solveTime := time.Now().UTC().Sub(chal.RevealAt)
    hours := int(solveTime.Hours())
    minutes := int(solveTime.Minutes()) % 60
    seconds := int(solveTime.Seconds()) % 60

    var solveTimeStr string
    if hours > 0 {
        solveTimeStr = fmt.Sprintf("%dh%dm%02d", hours, minutes, seconds)
    } else {
        solveTimeStr = fmt.Sprintf("%dm%02d", minutes, seconds)
    }

    blocks := []types.Block{
        {
            Type: "header",
            Text: &types.Text{
                Type: "plain_text",
                Text: ":sparkles: First Blood - " + chal.Name,
            },
        },
        {
			Type: "section",
			Fields: []types.Text{
				{
					Type: "mrkdwn",
					Text: "*Galadrimeur*:",
				},
				{
					Type: "mrkdwn",
					Text: "*Temps*:",
				},
				{
					Type: "plain_text",
					Text: solver.Name,
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("_%s_", solveTimeStr),
				},
			},
        },
    }

    message := types.Message{
        Blocks: blocks,
    }

	types.SendSlackWebhook(os.Getenv("SLACK_WEBHOOK_URI"), &message)
}
