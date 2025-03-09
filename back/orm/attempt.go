package orm

import (
	"math/rand/v2"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
const flagLength = 32

type Attempt struct {
	ID          uint64 `json:"id" gorm:"primaryKey" faker:"-"`
	UserID      uint64 `json:"user_id" gorm:"type:bigint"`
	ChallengeID uint64 `json:"challenge_id" gorm:"type:bigint"`
	Success     bool   `json:"-" faker:"-" gorm:"index"`
	Input       string `json:"-" faker:"word,lang=eng,len=16"`

	User      User      `json:"user" gorm:"foreignKey:UserID" faker:"-"`
	Challenge Challenge `json:"-" gorm:"foreignKey:ChallengeID" faker:"-"`
}

func GetSolvedAttempts(challengeId int) ([]Attempt, error) {
	var attempts []Attempt
	err := GormDB.
		Preload("User").
		Where("challenge_id = ? AND success = true", challengeId).
		Find(&attempts).Error
	if err != nil {
		return nil, err
	}
	return attempts, nil
}

func AsciiSum(s string) uint64 {
	var sum uint64
	for _, c := range s {
		sum += uint64(c)
	}
	return sum
}

func GenerateFlag(user *User, challengeName string) string {
	var flag strings.Builder
	rndSrc := rand.NewPCG(user.RandomSeed, AsciiSum(challengeName))
	flag.WriteString("GALA{")
	flag.Grow(flagLength)
	for i := 0; i < flagLength; i++ {
		flag.WriteByte(charset[rndSrc.Uint64()%uint64(len(charset))])
	}
	flag.WriteRune('}')
	return flag.String()
}

func VerifyFlag(user *User, challengeName, flag string) bool {
	return flag == GenerateFlag(user, challengeName)
}

func HasSolved(challengeId int, userId uint64) bool {
	var attempt *Attempt
	err := GormDB.
		Where("challenge_id = ? AND user_id = ? AND success = true", challengeId, userId).
		First(&attempt).Error
	return err == nil && attempt != nil
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
		maxChallenges := rand.IntN(3) + 1
		for i := 0; i < maxChallenges; i++ {
			var success bool
			if i == maxChallenges-1 {
				success = true
			}
			challenge := challenges[rand.IntN(len(challenges))]
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
