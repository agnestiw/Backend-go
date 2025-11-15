package test

import (
	"bytes"
	"encoding/json"
	"fmt" 
	"io"
	"mime/multipart"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Variabel global untuk ID File yang baru di-upload
var (
	testCreatedFileID   string
	testCreatedFileName = "test_image.png"
	testFileContent     = "ini-adalah-konten-gambar-dummy" // Konten dummy, tidak harus gambar asli
	testFileType        = "image/png"
)

// TestFile_1_Upload_Endpoint menguji POST /files/upload
func TestFile_1_Upload_Endpoint(t *testing.T) {

	// Skenario 1: Berhasil Meng-upload File (sebagai Admin)
	t.Run("Positive - Upload File Successfully (as Admin)", func(t *testing.T) {
		// Buat body request multipart
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		// 1. Buat form field "file" dengan header Content-Type manual
		// Ini adalah perbaikan untuk lolos validasi allowedTypes di service
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition",
			fmt.Sprintf(`form-data; name="file"; filename="%s"`, testCreatedFileName))
		h.Set("Content-Type", testFileType) // <-- Mengatur Content-Type ke "image/png"

		part, err := writer.CreatePart(h)
		assert.NoError(t, err)

		// Tulis konten file dummy
		_, err = part.Write([]byte(testFileContent))
		assert.NoError(t, err)

		// 2. (WAJIB) Tambahkan "target_user_id" karena kita admin
		// Kita set admin (testSeededUser) sebagai pemilik file
		err = writer.WriteField("target_user_id", testSeededUser.ID.Hex())
		assert.NoError(t, err)

		writer.Close() // Tutup writer untuk finalisasi body

		// Buat request
		req := httptest.NewRequest("POST", "/api/mg/files/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)

		// Cek respons dari handler UploadFile
		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode, "Status code seharusnya 201 Created")
		assert.Equal(t, true, respBody["success"])

		// Cek data file yg dikembalikan
		data, ok := respBody["data"].(map[string]interface{})
		assert.True(t, ok, "Key 'data' seharusnya ada")
		assert.Equal(t, testCreatedFileName, data["original_name"], "OriginalName harus sama")
		assert.Equal(t, testSeededUser.ID.Hex(), data["owner_id"], "OwnerID harus sama dengan target_user_id")

		// Simpan ID untuk tes berikutnya
		id, ok := data["id"].(string)
		assert.True(t, ok)
		assert.NotEmpty(t, id, "ID File yang baru dibuat tidak boleh kosong")
		testCreatedFileID = id // Simpan ke variabel global
	})

	// Skenario 2: Gagal karena Tidak Ada Token (Unauthorized)
	t.Run("Negative - Unauthorized (No Token)", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/mg/files/upload", nil)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	// Skenario 3: Gagal karena Tidak Ada File (Bad Request)
	t.Run("Negative - Bad Request (No File)", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/mg/files/upload", nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "Seharusnya 400 Bad Request jika tidak ada file")
	})

	// Skenario 4: Gagal karena Admin tidak menyertakan target_user_id
	t.Run("Negative - Bad Request (Admin No Target ID)", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", testCreatedFileName)
		part.Write([]byte(testFileContent))
		writer.Close()

		req := httptest.NewRequest("POST", "/api/mg/files/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+testAuthToken) // Token kita adalah admin
		// Tidak ada writer.WriteField("target_user_id", ...)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "Seharusnya 400 Bad Request jika admin tidak menyertakan target_user_id")
	})
}

// TestFile_2_GetMetadataByID_Endpoint menguji GET /files/:id (Metadata JSON)
func TestFile_2_GetMetadataByID_Endpoint(t *testing.T) {
	assert.NotEmpty(t, testCreatedFileID, "Test GetByID gagal: testCreatedFileID kosong (Upload mungkin gagal)")

	// Skenario 1: Berhasil Mendapatkan Metadata File by ID
	t.Run("Positive - Get File Metadata Successfully", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/files/"+testCreatedFileID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, true, respBody["success"])

		data, ok := respBody["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, testCreatedFileID, data["id"])
		assert.Equal(t, testCreatedFileName, data["original_name"])
	})

	// Skenario 2: Gagal karena ID Tidak Ditemukan (Not Found)
	t.Run("Negative - Not Found ID", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex() // ID acak
		req := httptest.NewRequest("GET", "/api/mg/files/"+nonExistentID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	// Skenario 3: Gagal karena Format ID Salah
	t.Run("Negative - Invalid ID Format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/files/id-invalid-format", nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		// Repo Anda akan gagal di helper.ToObjectID(), dan service akan return 500
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}

// TestFile_3_GetContentByID_Endpoint menguji GET /files/open/:id (Konten File)
func TestFile_3_GetContentByID_Endpoint(t *testing.T) {
	assert.NotEmpty(t, testCreatedFileID, "Test GetContent gagal: testCreatedFileID kosong")

	// Skenario 1: Berhasil Mendapatkan Konten File by ID
	t.Run("Positive - Get File Content Successfully", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/files/open/"+testCreatedFileID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Baca body (konten file)
		respBodyBytes, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		// Verifikasi konten
		// Set Content-Type di handler GetContentByID Anda agar ini lolos
		// c.Set("Content-Type", file.FileType)
		assert.Equal(t, testFileContent, string(respBodyBytes), "Konten file tidak sesuai")
	})

	// Skenario 2: Gagal karena ID Tidak Ditemukan (Not Found)
	t.Run("Negative - Not Found ID", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		req := httptest.NewRequest("GET", "/api/mg/files/open/"+nonExistentID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}

// TestFile_4_GetAllFiles_Endpoint menguji GET /files
func TestFile_4_GetAllFiles_Endpoint(t *testing.T) {
	// Skenario 1: Berhasil Mendapatkan Semua File
	t.Run("Positive - Get All Files", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/files", nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, true, respBody["success"])

		data, ok := respBody["data"].([]interface{})
		assert.True(t, ok)
		assert.True(t, len(data) >= 1, "Seharusnya ada minimal 1 file")
	})
}

// TestFile_5_DeleteFile_Endpoint menguji DELETE /files/:id (Hard Delete)
func TestFile_5_DeleteFile_Endpoint(t *testing.T) {
	assert.NotEmpty(t, testCreatedFileID, "Test Delete gagal: testCreatedFileID kosong")

	// Skenario 1: Berhasil Hard Delete
	t.Run("Positive - Hard Delete Successfully", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/mg/files/"+testCreatedFileID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, true, respBody["success"])
		assert.Equal(t, "File deleted successfully", respBody["message"])

		// Verifikasi juga bahwa file fisik telah dihapus
		// (Asumsi 'uploadPath' dari service)
		// fileLocation := filepath.Join("./uploads", "namafile-yg-disimpan-di-db")
		// _, err = os.Stat(fileLocation)
		// assert.True(t, os.IsNotExist(err), "File fisik seharusnya sudah terhapus")
	})

	// Skenario 2: Verifikasi Get Metadata Gagal Setelah di Hard-Delete
	t.Run("Negative - Get Metadata After Hard Delete", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/files/"+testCreatedFileID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode, "Seharusnya 404 Not Found setelah di-hard-delete")
	})

	// Skenario 3: Verifikasi Get Konten Gagal Setelah di Hard-Delete
	t.Run("Negative - Get Content After Hard Delete", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/files/open/"+testCreatedFileID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode, "Seharusnya 404 Not Found setelah di-hard-delete")
	})

	// Skenario 4: Gagal Menghapus Data yang Sudah Dihapus
	t.Run("Negative - Delete Already Deleted", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/mg/files/"+testCreatedFileID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		// Handler Anda akan gagal di FindFileByID (ErrNoDocuments)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode, "Seharusnya 404 Not Found saat menghapus data yg sudah dihapus")
	})
}