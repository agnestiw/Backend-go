package mongo

import (
	"latihan2/app/model"
	"latihan2/app/model/mongo"
	mongoRepo "latihan2/app/repository/mongo"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

// GetAllAlumni godoc
// @Summary Mendapatkan daftar alumni
// @Description Menampilkan semua data alumni dengan pagination, sorting, dan search
// @Tags Alumni
// @Accept json
// @Produce json
// @Param page query int false "Nomor halaman (default 1)"
// @Param limit query int false "Jumlah data per halaman (default 10)"
// @Param sortBy query string false "Kolom pengurutan (default: _id)"
// @Param order query string false "Urutan pengurutan (asc/desc)"
// @Param search query string false "Kata kunci pencarian"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mg/alumni [get]
func GetAllAlumni(c *fiber.Ctx) error {
	// Parsing query params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "_id")
	order := strings.ToLower(c.Query("order", "asc"))
	search := c.Query("search", "")
	offset := (page - 1) * limit

	// Panggil Repo
	data, err := mongoRepo.GetAlumniRepo(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	total, err := mongoRepo.CountAlumniRepo(search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"page":   page,
			"limit":  limit,
			"total":  total,
			"pages":  (total + limit - 1) / limit,
			"sortBy": sortBy,
			"order":  order,
			"search": search,
		},
	})
}

// GetAlumniByID godoc
// @Summary Mendapatkan data alumni berdasarkan ID
// @Description Menampilkan detail alumni berdasarkan ID-nya
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "ID Alumni"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mg/alumni/{id} [get]
func GetAlumniByID(c *fiber.Ctx) error {
	id := c.Params("id")

	data, err := mongoRepo.GetAlumniByID(id)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Alumni tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// CreateAlumni godoc
// @Summary Menambahkan data alumni baru
// @Description Membuat data alumni baru berdasarkan input pengguna
// @Tags Alumni
// @Accept json
// @Produce json
// @Param request body model.CreateAlumniRequest true "Data Alumni Baru"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mg/alumni [post]
func CreateAlumni(c *fiber.Ctx) error {
	// 1. Parse DTO Request (dari app/model/alumni.go)
	var req model.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 2. Ambil UserID dari token (ini PENTING)
	// Kita asumsikan middleware sudah menaruh string ObjectID
	userIDHex, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "UserID tidak ditemukan di token",
		})
	}
	userIDObj, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Format UserID di token tidak valid",
		})
	}

	// 3. Konversi DTO ke Model Database (app/model/mongo/alumni.go)
	newAlumni := &mongo.Alumni{
		UserID:     userIDObj, // ID dari token
		NIM:        req.NIM,
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		NoTelepon:  &req.NoTelepon, // Asumsi DTO string, model *string
		Alamat:     &req.Alamat,
	}

	// 4. Panggil Repo
	createdData, err := mongoRepo.CreateAlumni(newAlumni)
	if err != nil {
		// TODO: Handle duplikat NIM/Email
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    createdData,
	})
}

// UpdateAlumni godoc
// @Summary Mengupdate data alumni
// @Description Mengubah data alumni berdasarkan ID
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "ID Alumni"
// @Param request body model.UpdateAlumniRequest true "Data Alumni Terbaru"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mg/alumni/{id} [put]
func UpdateAlumni(c *fiber.Ctx) error {
	id := c.Params("id")

	// Parse DTO Request (dari app/model/alumni.go)
	var req model.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Panggil Repo
	updatedData, err := mongoRepo.UpdateAlumni(id, req)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Alumni tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    updatedData,
	})
}

// SoftDeleteAlumni godoc
// @Summary Menghapus (soft delete) data alumni
// @Description Menandai alumni sebagai dihapus tanpa menghapus permanen
// @Tags Alumni
// @Accept json
// @Produce json
// @Param id path string true "ID Alumni"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/mg/alumni/soft-delete/{id} [delete]
func SoftDeleteAlumni(c *fiber.Ctx) error {
	id := c.Params("id")

	err := mongoRepo.SoftDeleteAlumni(id)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Alumni tidak ditemukan atau sudah dihapus",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil di soft delete",
	})
}
