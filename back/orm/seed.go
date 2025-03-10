package orm

import (
	"log"
	"math/rand/v2"
	"os"
	"reflect"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/lib/pq"
)

var (
	eventStart = time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC)
)

// Idempotent seeding

func seedRealChallenges() {
	challenges := []Challenge{
		{
			Name:        "elite encryption",
			Difficulty:  0,
			Categories:  []string{"crypto"},
			Description: "After years of research, we've developed an unbreakable encryption algorithm. Only the most skilled cryptographers will be able to decode this.",
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "top_secret_data.txt",
					Size:  0,
					URL:   "/factories/elite_encryption",
				},
			},
		},
		{
			Name:        "one trick",
			Difficulty:  1,
			Categories:  []string{"crypto", "reverse"},
			Description: "Some say good security needs multiple layers. Others believe in minimalism. Here, we went all in on minimalism.",
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "uber-secure encrypted file reader",
					Size:  0,
					URL:   "/factories/one_trick",
				},
			},
		},
		{
			Name:        "super elite encryption",
			Difficulty:  0,
			Categories:  []string{"crypto"},
			Description: "This time we found a better encryption algorithm. There is no way someone could break this one.",
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "super_top_secret_data.txt",
					Size:  0,
					URL:   "/factories/super_elite_encryption",
				},
			},
		},
		{
			Name:        "quack",
			Difficulty:  2,
			Categories:  []string{"stegano"},
			Description: "I love this photo! I wonder who took it...",
			Attachments: []Attachment{
				{
					Type:  "file",
					Title: "quack.jpg",
					Size:  624247,
					URL:   "/factories/quack",
				},
				{
					Type:  "url",
					Title: "source",
					Size:  0,
					URL:   "https://www.pexels.com/photo/depth-of-field-photography-of-mallard-duck-on-body-of-water-660266/",
				},
			},
		},
	}

	// Create manually every challenges to skip already existing ones
	for i, challenge := range challenges {
		challenge.RevealAt = eventStart.Add(time.Duration(24*i) * time.Hour)
		GormDB.Create(&challenge)
	}
}

func seedUsers() {
	var userCount int64
	GormDB.Model(&User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("[-] users already seeded")
		return
	}
	log.Println("[#] seeding users")
	const MAX_USERS = 10
	users := FakeUsers(MAX_USERS)
	if err := GormDB.Create(&users).Error; err != nil {
		log.Fatalf("[-] error seeding users: %s", err.Error())
	}
	log.Printf("[+] seeded %d users!", len(users))
}

func seedAttempts() {
	var attemptCount int64
	GormDB.Model(&Attempt{}).Count(&attemptCount)
	if attemptCount > 0 {
		log.Println("[-] attempts already seeded")
		return
	}
	log.Println("[#] seeding attempts")
	attempts := FakeAttempts()
	if err := GormDB.Create(&attempts).Error; err != nil {
		log.Fatalf("[-] error seeding attempts: %s", err.Error())
	}
	log.Printf("[+] seeded %d attempts!", len(attempts))
}

func Seed() {
	seedRealChallenges()
	seedUsers()
	seedAttempts()

	os.Exit(0)
}

// Custom faker provider for categories
func customCategory() {
	validCategories := []string{"web", "crypto", "forensics", "misc", "network", "reverse", "web"}
	faker.AddProvider("ctfCategory", func(v reflect.Value) (interface{}, error) {
		numCategories := rand.IntN(3) + 1
		categories := make([]string, numCategories)
		for i := 0; i < numCategories; i++ {
			randIndex := rand.IntN(len(validCategories))
			categories[i] = validCategories[randIndex]
		}
		// hacky fix to assign a []string to a pq.StringArray (reflect will not allow this)
		return pq.StringArray(categories), nil
	})
}

func init() {
	customCategory()
}
