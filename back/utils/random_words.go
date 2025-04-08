package utils

import (
	"fmt"
	"math/rand/v2"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	nouns = []string{
		"county",
		"dinner",
		"chemistry",
		"guitar",
		"message",
		"connection",
		"virus",
		"weakness",
		"beer",
		"contribution",
		"obligation",
		"article",
		"driver",
		"preference",
		"establishment",
		"association",
		"wife",
		"cigarette",
		"drama",
		"video",
		"assignment",
		"attention",
		"two",
		"wedding",
		"vehicle",
		"candidate",
		"development",
		"method",
		"estate",
		"leadership",
		"boyfriend",
		"affair",
		"replacement",
		"ambition",
		"cell",
		"activity",
		"manufacturer",
		"finding",
		"editor",
		"diamond",
		"education",
		"attitude",
		"collection",
		"version",
		"blood",
		"tension",
		"hat",
		"mood",
		"engineering",
		"mode",
		"person",
		"conversation",
		"woman",
		"statement",
		"courage",
		"cancer",
		"topic",
		"inflation",
		"definition",
		"poet",
		"computer",
		"organization",
		"employer",
		"cookie",
		"food",
		"information",
		"quantity",
		"tradition",
		"player",
		"patience",
		"client",
		"employee",
		"response",
		"resolution",
		"meat",
		"opinion",
		"consequence",
		"hospital",
		"employment",
		"magazine",
		"initiative",
		"river",
		"feedback",
		"secretary",
		"cabinet",
		"hearing",
		"extent",
		"mall",
		"presence",
		"supermarket",
		"coffee",
		"assistance",
		"volume",
		"relation",
		"hair",
	}
	
	
	adjectives = []string{
		"sour",
		"energetic",
		"fluttering",
		"bashful",
		"gullible",
		"limping",
		"far",
		"squalid",
		"depressed",
		"special",
		"mature",
		"woozy",
		"actually",
		"practical",
		"pink",
		"anxious",
		"unhappy",
		"cruel",
		"watery",
		"determined",
		"silly",
		"thundering",
		"administrative",
		"greedy",
		"happy",
		"adorable",
		"agreeable",
		"wiry",
		"erratic",
		"right",
		"complete",
		"purring",
		"faint",
		"far-flung",
		"plausible",
		"lewd",
		"acidic",
		"heavenly",
		"redundant",
		"thoughtful",
		"learned",
		"supreme",
		"innocent",
		"automatic",
		"wide",
		"third",
		"guttural",
		"vague",
		"shocking",
		"narrow",
		"powerful",
		"careless",
		"hushed",
		"zesty",
		"outgoing",
		"minor",
		"shiny",
		"gray",
		"flat",
		"white",
		"smoggy",
		"gigantic",
		"idiotic",
		"physical",
		"acceptable",
		"adjoining",
		"stupid",
		"measly",
		"handsomely",
		"annoying",
		"dramatic",
		"curved",
		"unequaled",
		"pastoral",
		"wry",
		"shivering",
		"abrasive",
		"doubtful",
		"conscious",
		"lopsided",
		"combative",
		"like",
		"amazing",
		"questionable",
		"loving",
		"vivacious",
		"polite",
		"two",
		"tense",
		"sordid",
		"ancient",
		"abandoned",
		"famous",
		"flashy",
		"blue",
		"ambiguous",
		"used",
		"necessary",
		"utter",
	}

	musicExtensions = []string{
		".flac",
		".mp3",
		".wav",
		".ogg",
		".m4a",
	}

	photoExtensions = []string{
		".jpg",
		".png",
		".jpeg",
		".bmp",
	}

	documentExtensions = []string{
		".pdf",
		".doc",
		".docx",
		".xls",
		".xlsx",
		".ppt",
		".pptx",
	}

	videoExtensions = []string{
		".mp4",
		".mkv",
		".avi",
		".mov",
	}

	caser = cases.Title(language.English)
)

func RandomNoun(title bool) string {
	s := nouns[rand.IntN(len(nouns))]
	if title {
		s = caser.String(s)
	}
	return s
}

func RandomAdjective(title bool) string {
	s := adjectives[rand.IntN(len(adjectives))]
	if title {
		s = caser.String(s)
	}
	return s
}

func RandomNounAdj(title bool) string {
	return fmt.Sprintf("%s %s", RandomNoun(title), RandomAdjective(title))
}

func RandomMusic(title bool) string {
	return fmt.Sprintf("%s %s%s", RandomAdjective(title), RandomNoun(title), musicExtensions[rand.IntN(len(musicExtensions))])
}

func RandomPhoto(title bool) string {
	return fmt.Sprintf("%s %s%s", RandomAdjective(title), RandomNoun(title), photoExtensions[rand.IntN(len(photoExtensions))])
}

func RandomDocument(title bool) string {
	return fmt.Sprintf("%s %s%s", RandomAdjective(title), RandomNoun(title), documentExtensions[rand.IntN(len(documentExtensions))])
}

func RandomVideo(title bool) string {
	return fmt.Sprintf("%s %s%s", RandomAdjective(title), RandomNoun(title), videoExtensions[rand.IntN(len(videoExtensions))])
}