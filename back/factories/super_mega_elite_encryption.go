package factories

import (
	"strings"
	"unicode"

	"github.com/ggmolly/galactf/middlewares"
	"github.com/ggmolly/galactf/orm"
	"github.com/gofiber/fiber/v2"
)

var (
	emojiAlphabet = map[rune]string{
		'a': "🍎",
		'b': "🍌",
		'c': "🌰",
		'd': "🍩",
		'e': "🍆",
		'f': "🍟",
		'g': "🍇",
		'h': "🌺",
		'i': "🍦",
		'j': "🃏",
		'k': "🥝",
		'l': "🍋",
		'm': "🍈",
		'n': "🥜",
		'o': "🍊",
		'p': "🍍",
		'q': "🧀",
		'r': "🍓",
		's': "🍭",
		't': "🍵",
		'u': "🦄",
		'v': "🌋",
		'w': "🍉",
		'x': "❌",
		'y': "🍋",
		'z': "⚡",
		'0': "0️⃣",
		'1': "1️⃣",
		'2': "2️⃣",
		'3': "3️⃣",
		'4': "4️⃣",
		'5': "5️⃣",
		'6': "6️⃣",
		'7': "7️⃣",
		'8': "8️⃣",
		'9': "9️⃣",
	}
	upperCaseIndicator = "❗"
)

// Encode un texte en emojis
func textToEmoji(text string) string {
	var emojiText strings.Builder
	for _, char := range text {
		if word, found := emojiAlphabet[unicode.ToLower(char)]; found {
			emojiText.WriteString(word)
			if unicode.IsUpper(char) {
				emojiText.WriteString(upperCaseIndicator)
			}
		} else {
			emojiText.WriteRune(char)
		}
	}
	return emojiText.String()
}

func GenerateEmojiStegano(c *fiber.Ctx) error {
	user := middlewares.ReadUser(c)
	flag := orm.GenerateFlag(user, "emoji_stegano")
	return c.Status(fiber.StatusOK).SendString(textToEmoji(flag))
}
