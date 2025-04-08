package factories

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/ggmolly/galactf/middlewares"
	"github.com/ggmolly/galactf/orm"
	"github.com/ggmolly/galactf/utils"
	"github.com/gofiber/fiber/v2"
)

var (
	dirs = []string{
		"documents",
		"musics",
		"photos",
		"homeworks",
	}
	fileTypes = map[int]func(bool) string{
		0: utils.RandomDocument,
		1: utils.RandomMusic,
		2: utils.RandomPhoto,
		3: utils.RandomVideo,
	}
	fileSizeMin = 1024
	fileSizeMax = 2*1024*1024 // 2 MB
	maxSize = int64(30*1024*1024) // 30 MB
	mutex = &sync.Mutex{}
)

func unmount(mountPath string) error {
	return exec.Command("umount", mountPath).Run()
}

func clean(mountPath, imagePath string) error {
	if err := exec.Command("umount", mountPath).Run(); err != nil {
		return err
	}
	if err := os.Remove(imagePath); err != nil {
		return err
	}
	return os.RemoveAll(mountPath)
}

func GenerateByteBounty(c *fiber.Ctx) error {
	user := middlewares.ReadUser(c)
	outputImage := fmt.Sprintf("%d.img", user.ID)
	mountPath := fmt.Sprintf("/mnt/byte_bounty_%d", user.ID)
	
	mutex.Lock()
	defer mutex.Unlock()

	// Allocate image
	if err := exec.Command("fallocate", "-l", "50M", outputImage).Run(); err != nil {
		log.Println("[byte_bounty] failed to allocate file:", err)
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to generate image")
	}

	// Format to exfat
	if err := exec.Command("mkfs.exfat", outputImage).Run(); err != nil {
		log.Println("[byte_bounty] failed to format file:", err)
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to generate image")
	}

	// Mount image
	os.MkdirAll(mountPath, os.ModePerm)
	if err := exec.Command("mount", "-t", "exfat", "-o", "loop", outputImage, mountPath).Run(); err != nil {
		log.Println("[byte_bounty] failed to mount file:", err, "args:", outputImage, mountPath)
		os.Remove(outputImage)
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to generate image")
	}
	
	for _, dir := range dirs {
		os.MkdirAll(filepath.Join(mountPath, dir), os.ModePerm)
	}


	var seed [32]byte
	binary.LittleEndian.PutUint64(seed[:], uint64(user.RandomSeed))
	chachaReader := rand.NewChaCha8(seed)

	size := int64(0)
	flagIdx := 14 // worst case scenario: it's the penultimate file created (30//2 = 15)
	for {
		flagIdx--
		randomSize := int64(rand.IntN(fileSizeMax-fileSizeMin) + fileSizeMin)
		size += randomSize
		if size > maxSize {
			break
		}

		dir := rand.IntN(len(dirs))
		fileName := fileTypes[dir](true)
		fakeFile, err := os.Create(filepath.Join(mountPath, dirs[dir], fileName))
		if err != nil {
			log.Println("[byte_bounty] failed to create file:", err)
			return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to generate image")
		}

		if flagIdx != 0 {
			_, err = io.CopyN(fakeFile, chachaReader, randomSize)
			if err != nil {
				fakeFile.Close()
				log.Println("[byte_bounty] failed to write file:", err)
				clean(mountPath, outputImage)
				return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to generate image")
			}
		} else {
			flag := orm.GenerateFlag(user, "byte bounty")
			io.CopyN(fakeFile, chachaReader, randomSize/4+1)
			if _, err := fakeFile.WriteString(flag); err != nil {
				fakeFile.Close()
				log.Println("[byte_bounty] failed to write file:", err)
				clean(mountPath, outputImage)
				return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to generate image")
			} else {
				io.CopyN(fakeFile, chachaReader, randomSize/4+1)
			}
		}
		fakeFile.Close()
		// Gzip the file inline
		if err := exec.Command("gzip", filepath.Join(mountPath, dirs[dir], fileName)).Run(); err != nil {
			log.Println("[byte_bounty] failed to gzip file:", err)
			clean(mountPath, outputImage)
			return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to generate image")
		}
	}

	// Unmount image & clean up
	if err := clean(mountPath, outputImage); err != nil {
		log.Println("[byte_bounty] failed to clean file:", err)
		return utils.RestStatusFactory(c, fiber.StatusInternalServerError, "failed to generate image")
	}
	defer os.Remove(outputImage)

	c.Set("Content-Disposition", "attachment; filename=hdd_dump.img")
	return c.SendFile(outputImage)
}

