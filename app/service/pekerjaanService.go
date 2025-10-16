package service

import (
	"database/sql"
	"errors"
	"fmt"
	"latihan2/app/model"
	"latihan2/app/repository"
	"log"

	// "os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	// "github.com/golang-jwt/jwt/v5"
)

func GetPekerjaanByID(c *fiber.Ctx) error {
	log.Printf("[DEBUG] Locals: role=%v, user_id=%v", c.Locals("role"), c.Locals("user_id"))

	role := c.Locals("role").(string)
	if role == "user" {
		return c.Status(404).JSON(fiber.Map{"error": "Data tidak ada"})
	}

	// ambil id
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid pekerjaan ID"})
	}

	pekerjaan, err := repository.GetPekerjaanByIDRepo(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch pekerjaan"})
	}

	// if role != "admin" && pekerjaan.IsDelete {
	// 	return c.Status(404).JSON(fiber.Map{"error": "Data tidak ada"})
	// }

	return c.JSON(fiber.Map{"data": pekerjaan, "success": true})
}

func GetPekerjaanByAlumniID(c *fiber.Ctx) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID alumni tidak valid"})
	}

	pekerjaanList, err := repository.GetPekerjaanByAlumniID(alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data pekerjaan"})
	}

	return c.JSON(fiber.Map{"success": true, "data": pekerjaanList})
}

func CreatePekerjaan(c *fiber.Ctx) error {
	var req model.CreatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	// userID := c.Locals("user_id").(int) // ambil dari JWT

	newPekerjaan, err := repository.CreatePekerjaan(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menambah data pekerjaan. Pastikan alumni_id valid."})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    newPekerjaan,
		"message": "Pekerjaan berhasil ditambahkan",
	})
}

func UpdatePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID tidak valid"})
	}

	var req model.UpdatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	updatedPekerjaan, err := repository.UpdatePekerjaan(id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan untuk diupdate"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengupdate data pekerjaan"})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    updatedPekerjaan,
		"message": "Pekerjaan berhasil diupdate",
	})
}

func GetPekerjaanService(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")
	offset := (page - 1) * limit

	sortByWhitelist := map[string]bool{
		"id":                    true,
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
		"deleted_at":            true,
		"is_delete":             true,
		"delete_by":             true,
	}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	pekerjaan, err := repository.GetPekerjaanRepo(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to fetch alumni",
			"details": err.Error(),
			"query":   fmt.Sprintf("search=%s, sortBy=%s, order=%s, limit=%d, offset=%d, role=%s", search, sortBy, order, limit, offset),
		})
	}

	total, err := repository.CountPekerjaanRepo(search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count pekerjaan alumni"})
	}

	response := model.PekerjaanResponse{
		Data: pekerjaan,
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

func SoftDeletePekerjaan(c *fiber.Ctx) error {
	pekerjaanID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pekerjaan tidak valid"})
	}

	userID, okUserID := c.Locals("user_id").(int)
	role, okRole := c.Locals("role").(string)

	if !okUserID || !okRole {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Akses tidak sah, informasi pengguna tidak ditemukan"})
	}

	err = repository.SoftDeletePekerjaan(pekerjaanID, userID, role)
	if err != nil {
		log.Println("--- DEBUG: Terjadi error di repository ---", err)
		switch err {
		case repository.ErrPekerjaanNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data pekerjaan tidak ditemukan"})
		case repository.ErrForbidden:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Anda tidak diizinkan untuk menghapus data ini"})
		default:

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Terjadi kesalahan pada server"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan berhasil dihapus",
	})
}


func RestorePekerjaanService(c *fiber.Ctx) error {
	idParam := c.Params("id")
	pekerjaanID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID pekerjaan tidak valid",
		})
	}

	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	err = repository.RestorePekerjaan(pekerjaanID, userID, role)
	if err != nil {
		if errors.Is(err, repository.ErrPekerjaanNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Data pekerjaan tidak ditemukan atau belum dihapus",
			})
		}
		if errors.Is(err, repository.ErrForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Kamu tidak punya izin untuk mengembalikan data ini",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal melakukan restore pekerjaan",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data pekerjaan berhasil direstore",
	})
}


func GetTrashPekerjaanByIDService(c *fiber.Ctx) error {
	// Ambil ID pekerjaan dari parameter URL
	pekerjaanID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pekerjaan tidak valid"})
	}

	// Ambil informasi user dari middleware JWT
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	// Panggil repository untuk mendapatkan data trash
	pekerjaan, err := repository.GetTrashPekerjaanByID(pekerjaanID, userID, role)
	if err != nil {
		if err == sql.ErrNoRows {
			// Jika tidak ditemukan (atau tidak punya akses), kirim 404
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data pekerjaan di trash tidak ditemukan"})
		}
		// Error server lainnya
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data dari trash"})
	}

	// Kirim response sukses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    pekerjaan,
	})
}

func HardDeletePekerjaanService(c *fiber.Ctx) error {
	// Ambil ID pekerjaan dari parameter URL
	pekerjaanID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID pekerjaan tidak valid"})
	}

	// Ambil informasi user dari middleware JWT
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)

	// Panggil repository dengan informasi user
	err = repository.HardDeletePekerjaan(pekerjaanID, userID, role)
	if err != nil {
		if err == repository.ErrPekerjaanNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Data pekerjaan tidak ditemukan, belum di-soft delete, atau Anda tidak memiliki akses",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Terjadi kesalahan pada server"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan berhasil dihapus secara permanen",
	})
}