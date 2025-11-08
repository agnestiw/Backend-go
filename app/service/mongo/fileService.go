package mongo

import (
	"fmt"
	"log"
	"time"

	mongoModel "latihan2/app/model/mongo"
	mongoRepo "latihan2/app/repository/mongo"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const uploadPath = "./uploads"

func toFileResponse(file *mongoModel.File, ownerID primitive.ObjectID) *mongoModel.FileResponse {
	return &mongoModel.FileResponse{
		ID:           file.ID.Hex(),
		FileName:     file.FileName,
		OriginalName: file.OriginalName,
		FilePath:     file.FilePath,
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		UploadedAt:   file.UploadedAt,
		UploadedBy:   file.UploadedBy.Hex(),
		OwnerID:      ownerID.Hex(),
	}
}

// UploadFile godoc
// @Summary Upload file baru
// @Description Mengunggah file (gambar/pdf) ke server dan simpan metadata ke MongoDB.
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File yang akan diupload"
// @Param target_user_id formData string false "Hanya digunakan oleh admin untuk menentukan pemilik file"
// @Success 201 {object} map[string]interface{} "File uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/mg/files/upload [post]
func UploadFile(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "No file uploaded", "error": err.Error(),
		})
	}
	loggedInUserIDHex, okUserID := c.Locals("userID").(string)
	role, okRole := c.Locals("role").(string)
	if !okUserID || !okRole {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false, "message": "Informasi user tidak ditemukan di token",
		})
	}
	loggedInUserID, err := primitive.ObjectIDFromHex(loggedInUserIDHex)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false, "message": "Format UserID di token tidak valid",
		})
	}

	var ownerID primitive.ObjectID
	targetUserIDHex := c.FormValue("target_user_id")

	if role == "admin" {
		if targetUserIDHex == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Admin harus menyertakan target_user_id",
			})
		}
		targetObjID, err := primitive.ObjectIDFromHex(targetUserIDHex)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Format target_user_id tidak valid",
			})
		}
		ownerID = targetObjID
	} else if role == "user" {
		if targetUserIDHex != "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "User tidak diperbolehkan menentukan target_user_id",
			})
		}
		ownerID = loggedInUserID
	} else {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("Role '%s' tidak diizinkan melakukan upload", role),
		})
	}

	contentType := fileHeader.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/jpg":       true,
		"application/pdf": true,
	}
	if !allowedTypes[contentType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("File type '%s' not allowed", contentType),
		})
	}

	maxSize := int64(0)
	switch contentType {
	case "image/jpeg", "image/png", "image/jpg":
		maxSize = 1 * 1024 * 1024 // 1 MB
	case "application/pdf":
		maxSize = 2 * 1024 * 1024 // 2 MB
	}
	if fileHeader.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": fmt.Sprintf("File size exceeds limit for type %s (max %d MB)", contentType, maxSize/(1024*1024)),
		})
	}
	ext := filepath.Ext(fileHeader.Filename)
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(uploadPath, newFileName)

	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Failed to create upload directory", "error": err.Error(),
		})
	}
	if err := c.SaveFile(fileHeader, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Failed to save file to disk", "error": err.Error(),
		})
	}

	fileModel := &mongoModel.File{
		FileName:     newFileName,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileSize:     fileHeader.Size,
		FileType:     contentType,
		UploadedBy:   loggedInUserID,
		OwnerID:      ownerID,
		UploadedAt:   time.Now(),
	}

	if err := mongoRepo.CreateFile(fileModel); err != nil {
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Failed to save file metadata", "error": err.Error(),
		})
	}

	log.Printf("[DEBUG-UploadFile] file saved. ID: %s, OwnerID: %s", fileModel.ID.Hex(), ownerID.Hex())

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "File uploaded successfully",
		"data":    toFileResponse(fileModel, ownerID),
	})
}

// GetAllFiles godoc
// @Summary Ambil semua file
// @Description Mengambil daftar semua file yang tersimpan di database.
// @Tags Files
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Files retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/mg/files [get]
func GetAllFiles(c *fiber.Ctx) error {
	files, err := mongoRepo.FindAllFiles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Failed to get files", "error": err.Error(),
		})
	}

	var responses []mongoModel.FileResponse
	for _, file := range files {
		responses = append(responses, *toFileResponse(&file, file.OwnerID))
	}

	return c.JSON(fiber.Map{
		"success": true, "message": "Files retrieved successfully", "data": responses,
	})
}

// GetFileByID godoc
// @Summary Ambil file berdasarkan ID
// @Description Mengambil metadata file berdasarkan ID-nya.
// @Tags Files
// @Produce json
// @Security BearerAuth
// @Param id path string true "File ID"
// @Success 200 {object} map[string]interface{} "File retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Failure 500 {object} map[string]interface{} "Database error"
// @Router /api/mg/files/{id} [get]
func GetFileByID(c *fiber.Ctx) error {
	id := c.Params("id")
	file, err := mongoRepo.FindFileByID(id)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false, "message": "File not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Database error", "error": err.Error(),
		})
	}

	// ownerID = file.OwnerID
	return c.JSON(fiber.Map{
		"success": true, "message": "File retrieved successfully", "data": toFileResponse(file, file.OwnerID),
	})
}

func GetContentByID(c *fiber.Ctx) error {
	idHex := c.Params("id")
	fileID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid file ID format",
		})
	}

	loggedInUserIDHex, okUser := c.Locals("userID").(string)
	role, okRole := c.Locals("role").(string)
	if !okUser || !okRole {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid atau user belum login",
		})
	}

	loggedInUserID, err := primitive.ObjectIDFromHex(loggedInUserIDHex)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Format user ID tidak valid",
		})
	}

	file, err := mongoRepo.OpenFileByID(fileID)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "File tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil file dari database",
			"error":   err.Error(),
		})
	}

	if role == "user" && file.OwnerID != loggedInUserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Kamu tidak punya izin untuk mengakses file ini",
		})
	}

	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "File fisik tidak ditemukan di server",
		})
	}

	c.Set("Content-Type", file.FileType)
	c.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", file.OriginalName))
	return c.SendFile(file.FilePath)
}

// DeleteFile godoc
// @Summary Hapus file
// @Description Menghapus file dari server dan metadata dari database (hanya pemilik atau admin).
// @Tags Files
// @Produce json
// @Security BearerAuth
// @Param id path string true "File ID"
// @Success 200 {object} map[string]interface{} "File deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/mg/files/{id} [delete]
func DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")

	file, err := mongoRepo.FindFileByID(id)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false, "message": "File not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Database error", "error": err.Error(),
		})
	}

	if err := os.Remove(file.FilePath); err != nil {
		fmt.Println("Warning: Failed to delete file from storage:", err)
	}

	if err := mongoRepo.DeleteFile(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Failed to delete file metadata", "error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true, "message": "File deleted successfully",
	})
}

