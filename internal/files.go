package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/gofiber/fiber/v2"
)

// File represents a file stored in a namespace
type File struct {
	Filename string `json:"filename"`
	Data     []byte `json:"-"`    // don't expose raw bytes in API
	Hash     string `json:"hash"` // SHA256 hash
	FileSize int    `json:"size"`
}

// UploadFile handles file uploads and stores them in the EmberDB store
func UploadFile(c *fiber.Ctx) error {
	namespace := c.Params("namespace") // must provide namespace as query param
	key := c.Params("key")
	if namespace == "" {
		return fiber.NewError(fiber.StatusBadRequest, "namespace query parameter is required")
	}

	// Parse uploaded file
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

	newFile := File{
		Filename: fileHeader.Filename,
		Data:     fileBytes,
		Hash:     generateHash(fileBytes),
		FileSize: int(fileHeader.Size),
	}

	// Insert file into EmberDB store as Metadata
	store := &DataStore // global store

	store.Mu.Lock()
	defer store.Mu.Unlock()

	// Create namespace if it doesn't exist
	if store.Namespaces == nil {
		store.Namespaces = make(map[string]*Namespace)
	}
	ns, exists := store.Namespaces[namespace]
	if !exists {
		ns = &Namespace{
			Name: namespace,
			Data: make(map[string]Metadata),
		}
		store.Namespaces[namespace] = ns
	}

	// Check if key already exists
	if _, exists := ns.Data[key]; exists {
		return fiber.NewError(fiber.StatusConflict, "key already exists")
	}

	// Store file as Metadata
	ns.Data[key] = Metadata{
		Type:  TypeFile,
		Value: newFile,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "file uploaded successfully",
		"namespace": namespace,
		"key":       key,
		"hash":      newFile.Hash,
		"size":      newFile.FileSize,
	})
}

func generateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
