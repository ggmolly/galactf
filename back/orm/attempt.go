package orm

import "math/rand"

type Attempt struct {
	ID          uint64 `json:"id" gorm:"primaryKey" faker:"-"`
	UserID      uint64 `json:"user_id" gorm:"type:bigint"`
	ChallengeID uint64 `json:"challenge_id" gorm:"type:bigint"`
	Success     bool   `json:"success" faker:"-" gorm:"index"`
	Input       string `json:"input" faker:"word,lang=eng,len=16"`

	User      User      `json:"user" gorm:"foreignKey:UserID" faker:"-"`
	Challenge Challenge `json:"challenge" gorm:"foreignKey:ChallengeID" faker:"-"`
}

func FakeAttempts() []Attempt {
	var users []User
	if err := GormDB.Find(&users).Error; err != nil {
		panic(err)
	}

	var challenges []Challenge
	if err := GormDB.Find(&challenges).Error; err != nil {
		panic(err)
	}

	var attempts []Attempt
	for _, user := range users {
		maxChallenges := rand.Intn(3) + 1
		for i := 0; i < maxChallenges; i++ {
			var success bool
			if i == maxChallenges-1 {
				success = true
			}
			challenge := challenges[rand.Intn(len(challenges))]
			attempts = append(attempts, Attempt{
				UserID:      user.ID,
				ChallengeID: challenge.ID,
				Success:     success,
				Input:       "input",
			})
		}
	}

	return attempts
}
