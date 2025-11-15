package test

import (
	"bytes"
	"encoding/json"
	"latihan2/app/model/mongo"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var testCreatedPekerjaanID string

// TestPekerjaan_1_Create_Endpoint menguji POST /pekerjaan
func TestPekerjaan_1_Create_Endpoint(t *testing.T) {
	assert.NotEmpty(t, testSeededAlumniID, "Seeded Alumni ID tidak boleh kosong")
	seededAlumniObjectID, _ := primitive.ObjectIDFromHex(testSeededAlumniID)

	// Skenario 1: Berhasil Membuat Pekerjaan
	t.Run("Positive - Create Pekerjaan Successfully", func(t *testing.T) {
		reqBody := mongo.Pekerjaan{
			AlumniID:          seededAlumniObjectID,
			NamaPerusahaan:    "PT. Tester Jaya",
			PosisiJabatan:     "Quality Assurance",
			BidangIndustri:    "Teknologi",
			LokasiKerja:       "Jakarta",
			GajiRange:         "Rp 10jt - 15jt",
			TanggalMulaiKerja: time.Now().AddDate(-1, 0, 0),
			StatusPekerjaan:   "Full-time",
			Deskripsi:         "Melakukan pengujian pada aplikasi.",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/mg/pekerjaan", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode, "Status code seharusnya 201 Created")

		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, "PT. Tester Jaya", respBody["nama_perusahaan"])
		assert.Equal(t, testSeededAlumniID, respBody["alumni_id"], "AlumniID harus sama dengan yg di-seed")

		id, ok := respBody["id"].(string)
		assert.True(t, ok)
		assert.NotEmpty(t, id, "ID Pekerjaan yang baru dibuat tidak boleh kosong")
		testCreatedPekerjaanID = id
	})

	// Skenario 2: Gagal karena Tidak Ada Token (Unauthorized)
	t.Run("Negative - Unauthorized (No Token)", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/mg/pekerjaan", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	// Skenario 3: Gagal karena Body Request Tidak Valid (Bad Request)
	t.Run("Negative - Invalid Request Body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/mg/pekerjaan", bytes.NewBufferString(`{"nama_perusahaan": "abc",`)) // JSON rusak
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

// TestPekerjaan_2_GetByID_Endpoint menguji GET /pekerjaan/:id
func TestPekerjaan_2_GetByID_Endpoint(t *testing.T) {
	assert.NotEmpty(t, testCreatedPekerjaanID, "Test GetByID gagal: testCreatedPekerjaanID kosong (Create mungkin gagal)")

	// Skenario 1: Berhasil Mendapatkan Pekerjaan by ID
	t.Run("Positive - Get Pekerjaan Successfully", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/pekerjaan/"+testCreatedPekerjaanID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, testCreatedPekerjaanID, respBody["id"])
		assert.Equal(t, "PT. Tester Jaya", respBody["nama_perusahaan"])
	})

	// Skenario 2: Gagal karena ID Tidak Ditemukan (Not Found)
	t.Run("Negative - Not Found ID", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		req := httptest.NewRequest("GET", "/api/mg/pekerjaan/"+nonExistentID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	// Skenario 3: Gagal karena Format ID Salah
	t.Run("Negative - Invalid ID Format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/pekerjaan/ini-bukan-object-id", nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

// TestPekerjaan_3_GetByAlumniID_Endpoint menguji GET /pekerjaan/alumni/:alumni_id
func TestPekerjaan_3_GetByAlumniID_Endpoint(t *testing.T) {
	assert.NotEmpty(t, testSeededAlumniID, "Test GetByAlumniID gagal: testSeededAlumniID kosong")

	// Skenario 1: Berhasil Mendapatkan Pekerjaan by Alumni ID
	t.Run("Positive - Get Pekerjaan By Alumni ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/pekerjaan/alumni/"+testSeededAlumniID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		
		assert.Equal(t, testSeededAlumniID, respBody["alumni_id"])
		data, ok := respBody["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 1, "Seharusnya ada 1 pekerjaan untuk alumni ini")
	})

	// Skenario 2: Gagal karena Alumni ID Tidak Ditemukan (Not Found)
	t.Run("Negative - Alumni ID Not Found", func(t *testing.T) {
		nonExistentAlumniID := primitive.NewObjectID().Hex()
		req := httptest.NewRequest("GET", "/api/mg/pekerjaan/alumni/"+nonExistentAlumniID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode, "Seharusnya 404 jika tidak ada data")
	})
}

// TestPekerjaan_4_GetAll_Endpoint menguji GET /pekerjaan
func TestPekerjaan_4_GetAll_Endpoint(t *testing.T) {
	// Skenario 1: Berhasil Mendapatkan Semua (dengan search)
	t.Run("Positive - Get All With Search", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/pekerjaan?search=Tester%20Jaya", nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		meta, ok := respBody["meta"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(1), meta["total"], "Total data seharusnya 1")
		data, ok := respBody["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 1)
	})

	// Skenario 2: Gagal karena Tidak Ada Token (Unauthorized)
	t.Run("Negative - Unauthorized (No Token)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/pekerjaan", nil)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}

// TestPekerjaan_5_Update_Endpoint menguji PUT /pekerjaan/:id
func TestPekerjaan_5_Update_Endpoint(t *testing.T) {
	assert.NotEmpty(t, testCreatedPekerjaanID, "Test Update gagal: testCreatedPekerjaanID kosong")

	// Skenario 1: Berhasil Update Pekerjaan
	t.Run("Positive - Update Pekerjaan Successfully", func(t *testing.T) {
		reqBody := mongo.UpdatePekerjaanRequest{
			NamaPerusahaan: "PT. Tester Jaya (Updated)",
			PosisiJabatan:  "Senior QA",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PUT", "/api/mg/pekerjaan/"+testCreatedPekerjaanID, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		
		assert.Equal(t, "PT. Tester Jaya (Updated)", respBody["nama_perusahaan"])
		assert.Equal(t, "Senior QA", respBody["posisi_jabatan"])
	})

	// Skenario 2: Gagal karena ID Tidak Ditemukan (Not Found)
	t.Run("Negative - Not Found ID", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		reqBody := mongo.UpdatePekerjaanRequest{NamaPerusahaan: "Update Gagal"}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PUT", "/api/mg/pekerjaan/"+nonExistentID, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}

// TestPekerjaan_6_SoftDelete_Endpoint menguji DELETE /pekerjaan/soft-delete/:id
func TestPekerjaan_6_SoftDelete_Endpoint(t *testing.T) {
	assert.NotEmpty(t, testCreatedPekerjaanID, "Test SoftDelete gagal: testCreatedPekerjaanID kosong")

	// Skenario 1: Berhasil Soft Delete
	t.Run("Positive - Soft Delete Successfully", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/mg/pekerjaan/soft-delete/"+testCreatedPekerjaanID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, "Pekerjaan successfully soft deleted", respBody["message"])
	})

	// Skenario 2: Verifikasi Get By ID Gagal Setelah di Soft-Delete (Not Found)
	t.Run("Negative - Get After Soft Delete", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/pekerjaan/"+testCreatedPekerjaanID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode, "Seharusnya 404 Not Found setelah soft delete")
	})

	// Skenario 3: Gagal Menghapus Data yang Sudah Dihapus (Bad Request/Not Found)
	t.Run("Negative - Delete Already Deleted", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/mg/pekerjaan/soft-delete/"+testCreatedPekerjaanID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "Seharusnya 400 Bad Request saat menghapus data yg sudah di-soft-delete")
	})
}