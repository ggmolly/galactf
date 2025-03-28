package orm

import (
	"log"
	"math/rand/v2"
	"os"
	"reflect"
	"sort"
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
			Name:        "super elite encryption",
			Difficulty:  0,
			Categories:  []string{"stegano"},
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
			Name:        "more or less",
			Difficulty:  0,
			Categories:  []string{"misc"},
			Description: "Everyone knows what it is, but nobody knows how to write it.",
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "run_me.bf",
					Size:  0,
					URL:   "/factories/more_or_less",
				},
			},
		},
		{
			Name:        "cookie monster",
			Description: "Cookies are the most important part of the internet, everyone uses them, but a lot of developers don't know how to use them properly.",
			Difficulty:  1,
			Categories:  []string{"web", "pwn"},
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "secure webpage",
					Size:  0,
					URL:   "/factories/cookie_monster",
				},
			},
		},
		{
			Name:        "cookie monster²",
			Description: "Cookies are the most important part of the internet, everyone uses them, but a lot of developers don't know how to use them securely.",
			Difficulty:  3,
			Categories:  []string{"web", "pwn"},
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "secure webpage",
					Size:  0,
					URL:   "/factories/cookie_monster_squared",
				},
			},
		},
		{
			Name:        "claustrophobia",
			Description: "Ironically, overcoming claustrophobia often requires thinking outside of the box, even when the fear itself comes from feeling trapped inside one.",
			Difficulty:  2,
			Categories:  []string{"web", "pwn"},
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "File explorer",
					Size:  0,
					URL:   "/factories/claustrophobia",
				},
			},
		},
		{
			Name:        "quack",
			Difficulty:  1,
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
		{
			Name:        "bobby's library",
			Difficulty:  1,
			Categories:  []string{"web", "pwn"},
			Description: "Bobby made an online book library using sqlite! Check it out!",
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "bobby's website",
					Size:  0,
					URL:   "/factories/bobby_library",
				},
			},
		},
		{
			Name:        "cat image",
			Difficulty:  1,
			Categories:  []string{"stegano"},
			Description: "This is my cat, his name is Pixel! Isn't he adorable?",
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "image_of_a_cat.bmp",
					Size:  0,
					URL:   "/factories/cat_image",
				},
			},
		},
		{
			Name:        "(un)secure notes",
			Difficulty:  1,
			Categories:  []string{"web", "pwn"},
			Description: "the app where your secrets are never quite as safe as you think. Organize your thoughts, jot down your ideas, and try to keep them locked up…",
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "(un)secure notes",
					Size:  0,
					URL:   "/factories/unsecure_notes",
				},
			},
		},
		{
			Name:        "exclusive club",
			Difficulty:  3,
			Categories:  []string{"crypto", "reverse"},
			Description: "Some say good security needs multiple layers. Others believe in minimalism. Here, we went all in on minimalism.",
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "uber-secure encrypted file reader",
					Size:  0,
					URL:   "/factories/exclusive_club",
				},
			},
		},
		{
			Name:        "quiet riot code",
			Difficulty:  3,
			Categories:  []string{"stegano"},
			Description: "Quiet Riot Code (version 4) allows covert communication by embedding messages in ordinary digital activity, making them indistinguishable from routine data.",
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "intercepted_data_extraction.txt",
					Size:  0,
					URL:   "/factories/quiet_riot_code",
				},
			},
		},
		{
			Name:        "just a moment...",
			Description: "We developed a new ultra secure WAF anti-bot system. But some users were reporting a performance issue.. is V8 too slow?",
			Difficulty:  4,
			Categories:  []string{"web", "reverse", "crypto"},
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "captcha demo",
					Size:  0,
					URL:   "/factories/just_a_moment",
				},
			},
		},
		{
			Name:        "calculator",
			Description: "Our company worked on a calculator that runs server side for blazing fast performaces on our 64 core quantum computer. Feel free to try it out!",
			Difficulty:  4,
			Categories:  []string{"web", "pwn"},
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "calculator",
					Size:  0,
					URL:   "/factories/calculator",
				},
			},
		},
		{
			Name:        "twisted mersenne",
			Description: "Ever wondered why is pseudo-random number generator must NOT be used in cryptography? This challenge will teach you why!\nAlso, the system admin has weak credentials, but he instead heavily relies on his implementation of TOTP (rumored to be a simple random.randint).",
			Difficulty:  5,
			Categories:  []string{"crypto", "web", "pwn"},
			Attachments: []Attachment{
				{
					Type:  "url",
					Title: "milksad admin panel",
					Size:  0,
					URL:   "/factories/twisted_mersenne",
				},
			},
		},
	}

	// Sort by difficulty ascending
	sort.Slice(challenges, func(i, j int) bool {
		return challenges[i].Difficulty < challenges[j].Difficulty
	})

	// Create manually every challenges to skip already existing ones
	var daysOffset int
	for i := range challenges {
		challenges[i].RevealAt = eventStart.Add(time.Duration(24*(i+daysOffset)) * time.Hour)
		// If saturday or sunday, add two or one days
		if challenges[i].RevealAt.Weekday() == time.Saturday {
			challenges[i].RevealAt = challenges[i].RevealAt.Add(time.Hour * 48)
			daysOffset += 2
		}
		if challenges[i].RevealAt.Weekday() == time.Sunday {
			challenges[i].RevealAt = challenges[i].RevealAt.Add(time.Hour * 24)
			daysOffset += 1
		}
		// Closed day on the 21st of April
		if challenges[i].RevealAt.Day() == 21 {
			challenges[i].RevealAt = challenges[i].RevealAt.Add(time.Hour * 24)
			daysOffset += 1
		}
		GormDB.Create(&challenges[i])
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
