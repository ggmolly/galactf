package routes

import (
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/ggmolly/galactf/cache"
	"github.com/ggmolly/galactf/dto"
	"github.com/ggmolly/galactf/middlewares"
	"github.com/ggmolly/galactf/orm"
	protobuf "github.com/ggmolly/galactf/proto"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

const challengesCacheKey = "chal%d"
const challengesCacheTTL = time.Hour * 3

func GetChallenges(c *fiber.Ctx) error {
	user := middlewares.ReadUser(c)

	cachedChals, err := readCachedChallenges(user.ID)
	if err == nil {
		return utils.RestStatusFactoryData(c, fiber.StatusOK, cachedChals, "")
	}

	chals, err := orm.GetChallengeStats(user.ID)
	status := fiber.StatusOK
	if err != nil {
		log.Printf("[-] error fetching challenges: %s", err.Error())
		status = fiber.StatusInternalServerError
	}

	cache.WriteInterface(fmt.Sprintf(challengesCacheKey, user.ID), chals, challengesCacheTTL)
	return utils.RestStatusFactoryData(c, status, chals, "")
}

func GetChallenge(c *fiber.Ctx) error {
	user := middlewares.ReadUser(c)
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.RestStatusFactory(c, fiber.StatusBadRequest, "invalid id")
	}
	chal, err := orm.GetChallengeStatsById(id, user.ID)
	if err != nil {
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "error fetching challenge")
	}
	return utils.RestStatusFactoryData(c, fiber.StatusOK, chal, "")
}

func GetSolvers(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.RestStatusFactory(c, fiber.StatusBadRequest, "invalid id")
	}
	attempts, err := orm.GetSolvedAttempts(id)
	if err != nil {
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "error fetching challenge")
	}
	return utils.RestStatusFactoryData(c, fiber.StatusOK, attempts, "")
}

func SubmitFlag(c *fiber.Ctx) error {
	user := middlewares.ReadUser(c)
	chalId, err := c.ParamsInt("id")
	if err != nil {
		return utils.RestStatusFactory(c, fiber.StatusBadRequest, "Invalid challenge ID")
	}

	// Check if the user hasn't already solved the challenge
	if orm.HasSolved(chalId, user.ID) {
		return utils.RestStatusFactory(c, fiber.StatusConflict, "You've already solved this challenge")
	}

	dto, err := dto.ParseFlagSubmitDTO(c)
	if err != nil {
		return utils.RestStatusFactory(c, fiber.StatusBadRequest, "%s", err.Error())
	}
	chal, err := orm.GetChallengeById(chalId)
	if err != nil {
		return err
	}

	// Check if the challenge is locked'
	if chal.RevealAt.After(time.Now().UTC()) {
		return utils.RestStatusFactory(c, fiber.StatusForbidden, "This challenge is locked")
	}

	isValid := orm.VerifyFlag(user, chal.Name, dto.Flag)
	attempt := &orm.Attempt{
		UserID:      user.ID,
		ChallengeID: chal.ID,
		Success:     isValid,
		Input:       dto.Flag,
	}
	if err := orm.GormDB.Create(attempt).Error; err != nil {
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "Failed to submit flag")
	} else {
		orm.InvalidateLeaderboardCache()
	}

	var firstBlood bool
	var solveCount int64
	if err := orm.GormDB.
		Model(&orm.Attempt{}).
		Where("challenge_id = ? AND success = true", chal.ID).
		Count(&solveCount).
		Error; err == nil && solveCount == 1 { // silently ignore error
		firstBlood = true
	}

	// Broadcast the solve event to all connected clients
	event := protobuf.ChallengeAttempt{
		User: &protobuf.User{
			Id: user.ID,
		},
		ChallengeId: chal.ID,
		Success:     proto.Bool(isValid),
		FirstBlood:  proto.Bool(firstBlood),
	}

	if isValid {
		event.User.Name = &user.Name
	}

	// serialize the user name if their attempt is a first blood so the client
	// can display a toast message
	if firstBlood {
		orm.SendFirstBlood(chal, user)
	}

	Broadcast(protobuf.WS_CHALLENGE_ATTEMPT, &event)

	if !isValid { // Return an HTTP 201 (Created) status code if the flag is invalid but submitted successfully
		return utils.RestStatusFactory(c, fiber.StatusCreated, "Invalid flag! Try again.")
	} else { // Otherwise, return an HTTP 200 (OK) status code
		return utils.RestStatusFactory(c, fiber.StatusOK, "Congratulations, you've solved the challenge!")
	}
}

func readCachedChallenges(userID uint64) ([]orm.ChallengeStats, error) {
	b, err := cache.RedisDb.Get(cache.RedisCtx, fmt.Sprintf(challengesCacheKey, userID)).Bytes()
	if err == redis.Nil {
		return nil, redis.Nil
	} else if err != nil {
		log.Println("[!] failed to read challenges from cache:", err)
		return nil, err
	}
	var chals []orm.ChallengeStats
	err = sonic.ConfigFastest.Unmarshal(b, &chals)
	if err != nil {
		log.Println("[!] failed to unmarshal challenges:", err)
		return nil, err
	}
	return chals, nil
}
