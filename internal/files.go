package internal

import (
	"crypto/sha256"
	"emberdb/state"
	"encoding/hex"
	"fmt"
	"io"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var mu sync.RWMutex

type File struct {
	Filename string `json:"filename"`
	Data     []byte `json:"-"`
	Hash     string `json:"hash"`
	FileSize int    `json:"size"`
}

func UploadFile(c *fiber.Ctx) error {
	key := c.Params("key")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		fmt.Println("file is required")
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}
	if fileHeader.Size > 100<<20 {
		return fiber.NewError(fiber.StatusBadRequest, "file size exceeds 100MB limit")
	}

	file, err := fileHeader.Open()

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "cannot open file")
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "cannot read file")
	}
	newFile := File{
		Filename: fileHeader.Filename,
		Data:     fileBytes,
		Hash:     generateHash(fileBytes),
		FileSize: int(fileHeader.Size),
	}

	mu.Lock()
	state.DataStore[key] = newFile
	mu.Unlock()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "file uploaded successfully",
		"key":     key,
	})
}

func generateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
