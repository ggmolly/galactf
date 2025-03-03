package orm

import (
	"log"
	"math/rand/v2"
	"os"
	"reflect"

	"github.com/go-faker/faker/v4"
	"github.com/lib/pq"
)

// Idempotent seeding

func seedChallenges() {
	var challengeCount int64
	GormDB.Model(&Challenge{}).Count(&challengeCount)
	if challengeCount > 0 {
		log.Println("[-] challenges already seeded")
		return
	}
	log.Println("[#] seeding challenges")
	const MAX_CHALLENGES = 22
	chals := GetFakeChallenges(MAX_CHALLENGES)
	if err := GormDB.Create(&chals).Error; err != nil {
		log.Fatalf("[-] error seeding challenges: %s", err.Error())
	}
	log.Printf("[+] seeded %d challenges!", len(chals))
}

func Seed() {
	seedChallenges()

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
		return []pq.StringArray{pq.StringArray(categories)}, nil
	})
}

func init() {
	customCategory()
}
