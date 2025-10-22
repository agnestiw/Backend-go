package mongo

import (
	"fmt"
	mongoModel "latihan2/app/model/mongo"
	mongoRepo "latihan2/app/repository/mongo"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)


func GetAllPekerjaan(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "created_at")
	order := strings.ToLower(c.Query("order", "asc"))
	search := c.Query("search", "")
	offset := (page - 1) * limit

	sortWhitelist := map[string]bool{
		"_id":                   true,
		"user_id":               true,
		"alumni_id":             true,
		"nama_perusahaan":       true,
		"posisi_jabatan":        true,
		"bidang_industri":       true,
		"lokasi_kerja":          true,
		"gaji_range":            true,
		"tanggal_mulai_kerja":   true,
		"tanggal_selesai_kerja": true,
		"status_pekerjaan":      true,
		"created_at":            true,
		"updated_at":            true,
	}

	if !sortWhitelist[sortBy] {
		sortBy = "created_at"
	}

	if order != "desc" && order != "asc" {
		order = "asc"
	}

	data, err := mongoRepo.GetPekerjaanRepo(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	total, err := mongoRepo.CountPekerjaanRepo(search)
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

func GetPekerjaanByID(c *fiber.Ctx) error {
	id := c.Params("id")
	data, err := mongoRepo.GetPekerjaanByIDRepo(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pekerjaan not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(data)
}


func CreatePekerjaan(c *fiber.Ctx) error {
	var req mongoModel.Pekerjaan
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Simulate user/alumni ownership
	// userID := c.Locals("userID")
	// if uid, ok := userID.(string); ok {
	// 	if objID, err := primitive.ObjectIDFromHex(uid); err == nil {
	// 		req.UserID = objID
	// 	}
	// }

	req.CreatedAt = time.Now()
	req.IsDelete = false

	newPekerjaan, err := mongoRepo.CreatePekerjaan(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(newPekerjaan)
}

func GetPekerjaanByAlumniID(c *fiber.Ctx) error {
	alumniID := c.Params("alumni_id")

	data, err := mongoRepo.GetPekerjaanByAlumniID(alumniID)
	if err != nil {
		fmt.Println("DEBUG: Terjadi error saat ambil data dari repo:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// fmt.Printf("DEBUG: Jumlah data ditemukan: %d\n", len(data))
	// for i, d := range data {
	// 	fmt.Printf("DEBUG: Data[%d] = %+v\n", i, d)
	// }

	if len(data) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Tidak ada data pekerjaan untuk alumni ini",
			"debug": fiber.Map{
				"alumni_id": alumniID,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"alumni_id": alumniID,
		"data":      data,
	})
}


func UpdatePekerjaan(c *fiber.Ctx) error {
	id := c.Params("id")
	var req mongoModel.UpdatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	data, err := mongoRepo.UpdatePekerjaan(id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(data)
}


func SoftDeletePekerjaan(c *fiber.Ctx) error {
	fmt.Println("DEBUG: Masuk ke SoftDeletePekerjaan handler")
	id := c.Params("id")
	fmt.Println("DEBUG: ID diterima dari URL =", id)

	userID, ok := c.Locals("userID").(string)
	if !ok {
		fmt.Println("DEBUG: userID tidak ditemukan di Locals — menggunakan dummy untuk testing")
		userID = "dummy_user_id"
	}
	role, ok := c.Locals("role").(string)
	if !ok {
		fmt.Println("DEBUG: role tidak ditemukan di Locals — menggunakan dummy role admin")
		role = "admin"
	}

	if err := mongoRepo.SoftDeletePekerjaan(id, userID, role); err != nil {
		fmt.Println("DEBUG: Error saat SoftDeletePekerjaan:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan successfully soft deleted",
	})
}

func RestorePekerjaan(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("userID").(string)
	role := c.Locals("role").(string)

	if err := mongoRepo.RestorePekerjaan(id, userID, role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan successfully restored",
	})
}

// ✅ Hard delete
func HardDeletePekerjaan(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("userID").(string)
	role := c.Locals("role").(string)

	if err := mongoRepo.HardDeletePekerjaan(id, userID, role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Pekerjaan permanently deleted",
	})
}

