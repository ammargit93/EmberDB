package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/gofiber/fiber/v2"
)

type FileResponse struct {
	Filename string `json:"filename"`
	Hash     string `json:"hash"`
	FileSize int    `json:"size"`
}

func UploadFile(c *fiber.Ctx) error {
	namespace := c.Params("namespace")
	key := c.Params("key")

	if namespace == "" || key == "" {
		return fiber.NewError(fiber.StatusBadRequest, "namespace and key are required")
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
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

	// ✅ Convert to internal.Value immediately
	value := FileValue(fileBytes)

	// ✅ Insert through the DB API (not direct map mutation)
	ok, err := DataStore.Insert(namespace, key, value)
	if err != nil {
		return err
	}
	if !ok {
		return fiber.NewError(fiber.StatusConflict, "key already exists")
	}

	// API-only metadata
	resp := FileResponse{
		Filename: fileHeader.Filename,
		Hash:     generateHash(fileBytes),
		FileSize: int(fileHeader.Size),
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "file uploaded successfully",
		"namespace": namespace,
		"key":       key,
		"file":      resp,
	})
}

func generateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
