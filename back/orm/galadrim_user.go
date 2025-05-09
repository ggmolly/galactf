package orm

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bytedance/sonic"
	"github.com/ggmolly/galactf/cache"
	"github.com/gofiber/fiber/v2"
)

type GaladrimUser struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
}

const (
	ivSize            = 16
	ivSizeHex         = ivSize * 2
	cipheredEmailSize = 16
)

var (
	galadrimClient = &http.Client{
		Timeout: time.Second * 3,
	}
	cachedKeyBytes        []byte = nil
	ErrNotConnected              = errors.New("not connected to forest")
	ErrInvalidCookie             = errors.New("invalid cookie")
	ErrCookieReadFailed          = errors.New("failed to read cookie")
	ErrForestReqFailed           = errors.New("failed to resolve forest user")
	ErrForestNonSuccess          = errors.New("forest returned non-success status code")
	ErrForestDecodeFailed        = errors.New("failed to decode forest user")
	ErrInvalidEmail              = errors.New("invalid email")
)

func pkcs7Unpad(data []byte) []byte {
	paddingLen := int(data[len(data)-1])
	if paddingLen > aes.BlockSize || paddingLen == 0 {
		return data
	}
	return data[:len(data)-paddingLen]
}

func decipher(cipheredEmail []byte, initializationVector []byte) (string, error) {
	if cachedKeyBytes == nil {
		if b, err := hex.DecodeString(os.Getenv("GALADRIM_COOKIE_KEY")); err != nil {
			log.Fatal("[!] failed to decode GALADRIM_COOKIE_KEY:", err)
		} else {
			cachedKeyBytes = b
		}
	}
	block, err := aes.NewCipher(cachedKeyBytes)
	if err != nil {
		return "", err
	}
	decrypted := make([]byte, len(cipheredEmail))
	mode := cipher.NewCBCDecrypter(block, initializationVector)
	mode.CryptBlocks(decrypted, cipheredEmail)

	// ciphered bytes are padded to 64 bytes, thus we have to undo this padding
	return string(pkcs7Unpad(decrypted)), nil
}

func getGaladrimUser(email string) (*GaladrimUser, error) {
	// Check if email is ASCII only and contains an @
	if !utf8.ValidString(email) {
		log.Println("[!] invalid email:", email)
		return nil, ErrInvalidEmail
	}
	if !strings.Contains(email, "@") {
		log.Println("[!] no @ in email:", email)
		return nil, ErrInvalidEmail
	}

	req, _ := http.NewRequest("GET", "https://forest.galadrim.fr/profileInfos?email="+email, nil)
	req.Header.Add("User-Agent", "galactf-client/1.0")
	resp, err := galadrimClient.Do(req)
	if err != nil {
		log.Println("[!] failed to request forest user:", err)
		return nil, ErrForestReqFailed
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("[!] forest returned non-success status code:", resp.StatusCode, "email:", email)
		return nil, ErrForestNonSuccess
	}

	var user GaladrimUser
	err = sonic.ConfigFastest.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		log.Println("[!] failed to decode forest user:", err)
		return nil, ErrForestDecodeFailed
	}

	log.Println("[forest] resolved email:", email, "to:", user.Username)
	return &user, nil
}

func GetUserFromCookie(c *fiber.Ctx) (*User, error) {
	cookie := c.Cookies("galactf-cookie", c.Cookies("email-token", ""))
	if cookie == "" {
		log.Println("no cookie?")
		return nil, ErrNotConnected
	}

	cookie = strings.Replace(cookie, "%3A", ":", 1)

	// Decipher cookie
	initializationVector := make([]byte, ivSize)
	var cipheredEmail bytes.Buffer

	// expected length is 16+64, and +1 for the ':' delimiter
	// the ciphered email could be larger than 64 bytes, hency why we use a buffer
	// cipheredEmailSize is the minimum size of the ciphered email
	if len(cookie) < ivSizeHex+cipheredEmailSize+1 {
		log.Println("[!] invalid cookie length:", len(cookie))
		return nil, ErrInvalidCookie
	}

	cipheredEmail.Grow(cipheredEmailSize * 4)

	// Decode IV from hex
	{
		if iv, err := hex.DecodeString(cookie[:ivSizeHex]); err != nil {
			log.Println("[!] failed to decode IV:", err)
			return nil, ErrCookieReadFailed
		} else {
			copy(initializationVector, iv)
		}
	}

	// Decode email from hex
	{
		if email, err := hex.DecodeString(cookie[ivSizeHex+1:]); err != nil {
			log.Println("[!] failed to decode email:", err)
			return nil, ErrCookieReadFailed
		} else {
			cipheredEmail.Write(email)
		}
	}
	plaintextEmail, err := decipher(cipheredEmail.Bytes(), initializationVector)
	if err != nil {
		log.Println("[!] failed to decipher cookie:", err)
		return nil, ErrCookieReadFailed
	}

	// If the source of the cookie is not galactf, set a new cookie that expires in 30d
	if c.Cookies("galactf-cookie") == "" {
		c.Set("Set-Cookie", "galactf-cookie="+cookie+"; Max-Age=2592000; Path=/; HttpOnly; SameSite=Strict")
	}

	// Check if user is cached
	user, err := readCachedGaladrimUser(plaintextEmail)
	if err == nil {
		return user, nil
	}

	galaUser, err := getGaladrimUser(plaintextEmail)
	if err != nil {
		log.Println("[?] error in galadrim user")
		return nil, err
	}

	// Resolve Galadrim user <=> galactf user with email
	var galactfUser User
	err = GormDB.Where("email = ?", plaintextEmail).First(&galactfUser).Error

	// If the user exists, cache it, and return it
	if err == nil {
		cache.WriteInterface(plaintextEmail, RedisUser{
			ID:         galactfUser.ID,
			Name:       galactfUser.Name,
			RandomSeed: galactfUser.RandomSeed,
		}, cache.GalaUserCacheTTL)
		return &galactfUser, nil
	}

	// Otherwise, create the user and cache it
	galactfUser.Email = plaintextEmail
	galactfUser.Name = galaUser.FullName
	if galactfUser.RandomSeed, err = GenerateRandomSeed(); err != nil {
		log.Println("[!] failed to generate random seed:", err)
		return nil, err
	}
	if err := GormDB.Create(&galactfUser).Error; err != nil {
		log.Println("[!] failed to create user:", err)
		return nil, err
	}

	cache.WriteInterface(plaintextEmail, RedisUser{
		ID:         galactfUser.ID,
		Name:       galactfUser.Name,
		RandomSeed: galactfUser.RandomSeed,
	}, cache.GalaUserCacheTTL)

	return user, nil
}

func readCachedGaladrimUser(email string) (*User, error) {
	r, err := cache.ReadCached[RedisUser](email)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:         r.ID,
		Name:       r.Name,
		RandomSeed: r.RandomSeed,
	}, err
}
