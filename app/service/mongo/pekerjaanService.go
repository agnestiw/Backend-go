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

// GetAllPekerjaan godoc
// @Summary      Mendapatkan semua data pekerjaan
// @Description  Mengambil daftar pekerjaan dengan pagination, sorting, dan search
// @Tags         Pekerjaan
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "Halaman"
// @Param        limit    query     int     false  "Jumlah data per halaman"
// @Param        sortBy   query     string  false  "Kolom pengurutan (default: created_at)"
// @Param        order    query     string  false  "Arah pengurutan (asc/desc)"
// @Param        search   query     string  false  "Kata kunci pencarian"
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/mg/pekerjaan [get]
// @Security     BearerAuth
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

// GetPekerjaanByID godoc
// @Summary      Mendapatkan pekerjaan berdasarkan ID
// @Description  Mengambil satu data pekerjaan berdasarkan ID
// @Tags         Pekerjaan
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID pekerjaan"
// @Success      200 {object} mongoModel.Pekerjaan
// @Failure      404 {object} map[string]interface{}
// @Router       /api/mg/pekerjaan/{id} [get]
// @Security     BearerAuth
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

// GetPekerjaanByAlumniID godoc
// @Summary      Get pekerjaan berdasarkan alumni ID
// @Description  Mengambil semua data pekerjaan yang dimiliki oleh seorang alumni berdasarkan ID alumni
// @Tags         Pekerjaan
// @Accept       json
// @Produce      json
// @Param        alumni_id   path      string  true  "ID Alumni"
// @Success      200  {object}  map[string]interface{}  "Daftar pekerjaan berdasarkan alumni"
// @Failure      404  {object}  map[string]string  "Tidak ada data pekerjaan untuk alumni ini"
// @Failure      500  {object}  map[string]string  "Terjadi kesalahan server"
// @Router       /api/mg/pekerjaan/alumni/{alumni_id} [get]
// @Security     BearerAuth
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

// CreatePekerjaan godoc
// @Summary      Membuat data pekerjaan baru
// @Description  Menambahkan data pekerjaan baru ke dalam database
// @Tags         Pekerjaan
// @Accept       json
// @Produce      json
// @Param        request body mongoModel.Pekerjaan true "Data pekerjaan baru"
// @Success      201 {object} mongoModel.Pekerjaan
// @Failure      400 {object} map[string]interface{}
// @Failure      500 {object} map[string]interface{}
// @Router       /api/mg/pekerjaan [post]
// @Security     BearerAuth
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

// UpdatePekerjaan godoc
// @Summary      Memperbarui data pekerjaan
// @Description  Mengupdate data pekerjaan berdasarkan ID
// @Tags         Pekerjaan
// @Accept       json
// @Produce      json
// @Param        id path string true "ID pekerjaan"
// @Param        request body mongoModel.UpdatePekerjaanRequest true "Data pekerjaan yang diperbarui"
// @Success      200 {object} mongoModel.Pekerjaan
// @Failure      400 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /api/mg/pekerjaan/{id} [put]
// @Security     BearerAuth
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

// SoftDeletePekerjaan godoc
// @Summary      Soft delete pekerjaan
// @Description  Menghapus sementara data pekerjaan (tidak benar-benar dihapus)
// @Tags         Pekerjaan
// @Accept       json
// @Produce      json
// @Param        id path string true "ID pekerjaan"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /api/mg/pekerjaan/soft-delete/{id} [delete]
// @Security     BearerAuth
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

// RestorePekerjaan godoc
// @Summary      Mengembalikan data pekerjaan yang terhapus
// @Description  Restore data pekerjaan yang dihapus secara soft delete
// @Tags         Pekerjaan
// @Accept       json
// @Produce      json
// @Param        id path string true "ID pekerjaan"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /api/mg/pekerjaan/restore/{id} [post]
// @Security     BearerAuth
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

// HardDeletePekerjaan godoc
// @Summary      Menghapus permanen data pekerjaan
// @Description  Hard delete data pekerjaan dari database
// @Tags         Pekerjaan
// @Accept       json
// @Produce      json
// @Param        id path string true "ID pekerjaan"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Router       /api/mg/pekerjaan/hard-delete/{id} [delete]
// @Security     BearerAuth
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

// GetTrashPekerjaanByID godoc
// @Summary      Mendapatkan pekerjaan yang dihapus berdasarkan ID
// @Description  Menampilkan data pekerjaan yang sudah dihapus (soft delete) berdasarkan ID tertentu
// @Tags         Pekerjaan
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID pekerjaan"
// @Success      200  {object}  map[string]interface{}  "Data pekerjaan yang sudah dihapus berhasil diambil"
// @Failure      400  {object}  map[string]interface{}  "Format ID tidak valid"
// @Failure      404  {object}  map[string]interface{}  "Data pekerjaan tidak ditemukan di trash"
// @Failure      500  {object}  map[string]interface{}  "Terjadi kesalahan pada server"
// @Security     BearerAuth
// @Router       /api/mg/pekerjaan/trash/{id} [get]
func GetTrashPekerjaan(c *fiber.Ctx) error {
	id := c.Params("id")
	role, _ := c.Locals("role").(string)

	pekerjaanList, err := mongoRepo.GetTrashPekerjaan(id, role)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": pekerjaanList,
	})
}
