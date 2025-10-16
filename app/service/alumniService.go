package service

import (
	"database/sql"
	"fmt"
	"latihan2/app/model"
	"latihan2/app/repository"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// func GetAllAlumni(c *fiber.Ctx) error {
// 	alumni, err := repository.GetAllAlumni()
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"success": false,
// 			"error":   "Gagal mengambil data alumni: " + err.Error(),
// 		})
// 	}
// 	return c.JSON(fiber.Map{
// 		"success": true,
// 		"data":    alumni,
// 		"message": "Data alumni berhasil diambil",
// 	})
// }

func GetAlumniByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	role := c.Locals("role").(string)

	alumni, err := repository.GetAlumniByID(id, role)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Alumni tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data alumni"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    alumni,
		"message": "Data alumni berhasil diambil",
	})
}

func CreateAlumni(c *fiber.Ctx) error {
	var req model.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	if req.NIM == "" || req.Nama == "" || req.Jurusan == "" || req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Field wajib (NIM, Nama, Jurusan, Email) harus diisi"})
	}

	// Ambil userID dari c.Locals (diasumsikan middleware menyimpannya di sini)
	// userID, ok := c.Locals("user_id").(int)
	// if !ok {
	// 	// Penanganan jika userID tidak ditemukan atau tipe datanya salah
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Akses tidak sah"})
	// }

	// Teruskan userID saat memanggil repository
	newAlumni, err := repository.CreateAlumni(req)
	if err != nil {
		// Jika error masih terjadi, kemungkinan karena constraint UNIQUE di database
		log.Println("Error creating alumni:", err) // Tambahkan log untuk debugging
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menambah alumni. Pastikan NIM dan email belum digunakan"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    newAlumni,
		"message": "Alumni berhasil ditambahkan",
	})
}

func UpdateAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var req model.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	updatedAlumni, err := repository.UpdateAlumni(id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Alumni tidak ditemukan untuk diupdate"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengupdate alumni"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    updatedAlumni,
		"message": "Alumni berhasil diupdate",
	})
}

func DeleteAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	err = repository.DeleteAlumni(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Alumni tidak ditemukan untuk dihapus"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus alumni"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil dihapus",
	})
}

func GetAlumniByTahunLulus(c *fiber.Ctx) error {
	tahunParam := c.Params("tahun")
	tahun, err := strconv.Atoi(tahunParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Tahun lulus tidak valid",
		})
	}

	data, total, err := repository.GetAlumniByTahunLulus(tahun)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal mengambil data alumni: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":                     true,
		"data":                        data,
		"total_alumni_gaji_lebih_4jt": total,
		"message":                     "Data alumni berdasarkan tahun lulus",
	})
}

func GetAlumniService(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")
	offset := (page - 1) * limit

	role := c.Locals("role").(string)

	alumni, err := repository.GetAlumniRepo(search, sortBy, order, limit, offset, role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to fetch alumni",
			"details": err.Error(),
			"query":   fmt.Sprintf("search=%s, sortBy=%s, order=%s, limit=%d, offset=%d, role=%s", search, sortBy, order, limit, offset, role),
		})
	}

	total, err := repository.CountAlumniRepo(search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count alumni"})
	}

	response := model.AlumniResponse{
		Data: alumni,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}

	return c.JSON(response)
}

func SoftDeleteAlumniService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	// hanya admin boleh delete
	role := c.Locals("role").(string)
	if role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Hanya admin yang bisa menghapus alumni"})
	}

	err = repository.SoftDeleteAlumniRepo(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Alumni tidak ditemukan atau sudah dihapus"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus alumni"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil dihapus (soft delete)",
	})
}
