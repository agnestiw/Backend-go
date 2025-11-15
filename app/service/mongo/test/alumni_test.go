package test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"latihan2/middleware" // <-- Impor middleware
	"latihan2/app/model"
	mongoModel "latihan2/app/model/mongo"
	mongoService "latihan2/app/service/mongo" // Impor service mongo
	"latihan2/database"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Variabel global untuk state test suite
var (
	testApp        *fiber.App
	testAuthToken  string // Token JWT untuk autentikasi
	testSeededUser mongoModel.User

	// Variabel untuk tes Alumni
	testCreatedAlumniID string // ID dari alumni yang baru dibuat (oleh tes alumni)
	testAlumniNIM       = "TEST-NIM-001" // NIM untuk tes create alumni

	// Variabel untuk tes Pekerjaan (diisi oleh TestMain)
	testSeededAlumniID     string // <-- TAMBAHAN: Variabel dari pekerjaan_test.go akan diisi di sini
	testSeededAlumniForPekerjaan mongoModel.Alumni // <-- TAMBAHAN
	testAlumniNIMForPekerjaan = "TEST-NIM-FOR-JOB-002" // <-- TAMBAHAN

	// Variabel login
	testSeededUsername = "alumni_tester_mongo"
	testSeededPassword = "alumni_pass123"
)

// TestMain akan dieksekusi sekali sebelum & sesudah semua tes di package ini.
func TestMain(m *testing.M) {
	// 1. SETUP
	err := godotenv.Load("../../../../.env")
	if err != nil {
		log.Fatalf("Gagal memuat file .env dari root: %v", err)
	}
	os.Setenv("JWT_SECRET_KEY", "16824af3-6b8e-4c3d-9f1e-2c4b5e6f7g8h") // Pastikan sama

	// Hubungkan ke DB
	database.InitMongoDB()
	if database.MongoDB == nil {
		log.Fatal("Koneksi MongoDB nil setelah InitMongoDB")
	}

	// Siapkan collections
	collUser := database.MongoDB.Collection("user")
	collAlumni := database.MongoDB.Collection("alumni")
	collPekerjaan := database.MongoDB.Collection("pekerjaan") // <-- TAMBAHAN

	// Bersihkan data lama
	_, _ = collUser.DeleteMany(context.Background(), bson.M{"username": testSeededUsername})
	// <-- TAMBAHAN: Bersihkan kedua NIM
	_, _ = collAlumni.DeleteMany(context.Background(), bson.M{"nim": bson.M{"$in": []string{testAlumniNIM, testAlumniNIMForPekerjaan}}})


	// 2. SEED USER (untuk mendapatkan token)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testSeededPassword), bcrypt.DefaultCost)
	testSeededUser = mongoModel.User{
		ID:        primitive.NewObjectID(),
		Username:  testSeededUsername,
		Email:     testSeededUsername + "@example.com",
		Password:  string(hashedPassword),
		Role:      "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = collUser.InsertOne(context.Background(), testSeededUser)
	if err != nil {
		log.Fatalf("Gagal seeding data user tes: %v", err)
	}

	// 3. SEED ALUMNI (UNTUK KEBUTUHAN TES PEKERJAAN)
	// <-- TAMBAHAN: Blok ini penting
	testSeededAlumniForPekerjaan = mongoModel.Alumni{
		ID:         primitive.NewObjectID(),
		UserID:     testSeededUser.ID,
		NIM:        testAlumniNIMForPekerjaan,
		Nama:       "Alumni untuk Tes Pekerjaan",
		Jurusan:    "Teknik Mesin",
		Angkatan:   2019,
		TahunLulus: 2023,
		Email:      "job-tester@alumni.com",
		CreatedAt:  time.Now(),
	}
	_, err = collAlumni.InsertOne(context.Background(), testSeededAlumniForPekerjaan)
	if err != nil {
		log.Fatalf("Gagal seeding data alumni tes: %v", err)
	}
	// Ini akan mengisi variabel global yg dibutuhkan oleh pekerjaan_test.go
	testSeededAlumniID = testSeededAlumniForPekerjaan.ID.Hex() 
	
	// Bersihkan pekerjaan yg mungkin nyangkut dari alumni ini
	_, _ = collPekerjaan.DeleteMany(context.Background(), bson.M{"alumni_id": testSeededAlumniForPekerjaan.ID})


	// 4. SETUP FIBER APP & ROUTES
	testApp = fiber.New()
	api := testApp.Group("/api/mg")

	// Rute Publik (Login)
	api.Post("/login", mongoService.LoginMongo)

	// Rute Terproteksi (Middleware)
	protectedm := api.Group("", middleware.AuthRequiredMongo()) 

	// Rute Alumni
	alumnim := protectedm.Group("/alumni")
	alumnim.Get("/", mongoService.GetAllAlumni)
	alumnim.Get("/:id/", mongoService.GetAlumniByID)
	alumnim.Post("/", mongoService.CreateAlumni)
	alumnim.Put("/:id", mongoService.UpdateAlumni)
	alumnim.Delete("/soft-delete/:id", mongoService.SoftDeleteAlumni)

	// Rute Pekerjaan
	pekerjaanm := protectedm.Group("/pekerjaan")
	pekerjaanm.Get("/", mongoService.GetAllPekerjaan)
	pekerjaanm.Get("/:id", mongoService.GetPekerjaanByID)
	pekerjaanm.Get("/alumni/:alumni_id", mongoService.GetPekerjaanByAlumniID)
	pekerjaanm.Post("/", mongoService.CreatePekerjaan)
	pekerjaanm.Put("/:id", mongoService.UpdatePekerjaan)
	pekerjaanm.Delete("/soft-delete/:id", mongoService.SoftDeletePekerjaan)

	// Rute File
	filem := protectedm.Group("/files")
	filem.Post("/upload", mongoService.UploadFile)
	filem.Get("/", mongoService.GetAllFiles)
	filem.Get("/:id", mongoService.GetFileByID)
	filem.Get("/open/:id", mongoService.GetContentByID)
	filem.Delete("/:id", mongoService.DeleteFile)

	// 5. LOGIN UNTUK MENDAPATKAN TOKEN
	loginBody := model.LoginRequest{Username: testSeededUsername, Password: testSeededPassword}
	bodyBytes, _ := json.Marshal(loginBody)
	req := httptest.NewRequest("POST", "/api/mg/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := testApp.Test(req, -1)
	if err != nil || resp.StatusCode != fiber.StatusOK {
		log.Fatalf("Gagal login saat setup test: %v (Status: %d)", err, resp.StatusCode)
	}

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	data, ok := respBody["data"].(map[string]interface{})
	if !ok {
		log.Fatal("Key 'data' tidak ditemukan di respons login")
	}
	testAuthToken, ok = data["token"].(string)
	if !ok || testAuthToken == "" {
		log.Fatal("Gagal mendapatkan token JWT dari respons login")
	}

	// 6. JALANKAN SEMUA TES
	exitCode := m.Run()

	// 7. TEARDOWN
	// Bersihkan data yang dibuat selama tes
	_, _ = collUser.DeleteMany(context.Background(), bson.M{"username": testSeededUsername})
	_, _ = collAlumni.DeleteMany(context.Background(), bson.M{"user_id": testSeededUser.ID})
	// <-- TAMBAHAN: Bersihkan pekerjaan yg dibuat
	_, _ = collPekerjaan.DeleteMany(context.Background(), bson.M{"alumni_id": testSeededAlumniForPekerjaan.ID})
	
	if database.MongoClient != nil {
		database.MongoClient.Disconnect(context.Background())
	}

	os.Exit(exitCode)
}

// =========================================================================
// SEMUA TES ALUMNI DI BAWAH INI TETAP SAMA, TIDAK PERLU DIUBAH
// =========================================================================

// TestAlumni_1_Create_Endpoint menguji endpoint POST /alumni
func TestAlumni_1_Create_Endpoint(t *testing.T) {
	// Skenario 1: Berhasil Membuat Alumni
	t.Run("Positive - Create Alumni Successfully", func(t *testing.T) {
		reqBody := model.CreateAlumniRequest{
			NIM:        testAlumniNIM,
			Nama:       "Nama Tester",
			Jurusan:    "Teknik Informatika",
			Angkatan:   2020,
			TahunLulus: 2024,
			Email:      "tester@alumni.com",
			NoTelepon:  "08123456789",
			Alamat:     "Jalan Uji Coba",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/api/mg/alumni", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testAuthToken) // <-- Set token

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode, "Status code seharusnya 201 Created")

		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)

		data, ok := respBody["data"].(map[string]interface{})
		assert.True(t, ok, "Key 'data' seharusnya ada")
		assert.Equal(t, testAlumniNIM, data["nim"])
		assert.Equal(t, testSeededUser.ID.Hex(), data["user_id"], "UserID harus sama dengan user yg login")

		// Simpan ID untuk tes berikutnya
		testCreatedAlumniID, ok = data["id"].(string)
		assert.True(t, ok)
		assert.NotEmpty(t, testCreatedAlumniID, "ID Alumni yang baru dibuat tidak boleh kosong")
	})

	// Skenario 2: Gagal karena Tidak Ada Token (Unauthorized)
	t.Run("Negative - Unauthorized (No Token)", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/mg/alumni", nil)
		req.Header.Set("Content-Type", "application/json")
		// Tidak ada header Authorization

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	// Skenario 3: Gagal karena Body Request Tidak Valid (Bad Request)
	t.Run("Negative - Invalid Request Body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/mg/alumni", bytes.NewBufferString(`{"nim": "123",`)) // JSON rusak
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

// TestAlumni_2_GetByID_Endpoint menguji GET /alumni/:id
func TestAlumni_2_GetByID_Endpoint(t *testing.T) {
	assert.NotEmpty(t, testCreatedAlumniID, "Test GetByID gagal: testCreatedAlumniID kosong (Create mungkin gagal)")

	// Skenario 1: Berhasil Mendapatkan Alumni by ID
	t.Run("Positive - Get Alumni Successfully", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/alumni/"+testCreatedAlumniID+"/", nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		data, ok := respBody["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, testCreatedAlumniID, data["id"])
		assert.Equal(t, testAlumniNIM, data["nim"])
	})

	// Skenario 2: Gagal karena ID Tidak Ditemukan (Not Found)
	t.Run("Negative - Not Found ID", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex() // ID acak yang pasti tidak ada
		req := httptest.NewRequest("GET", "/api/mg/alumni/"+nonExistentID+"/", nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	// Skenario 3: Gagal karena Format ID Salah (Server Error / Bad Request)
	t.Run("Negative - Invalid ID Format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/alumni/ini-bukan-object-id/", nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		// Service Anda mengembalikan 500 jika format ID salah di repo
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	// Skenario 4: Gagal karena Tidak Ada Token (Unauthorized)
	t.Run("Negative - Unauthorized (No Token)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/alumni/"+testCreatedAlumniID+"/", nil)
		// Tidak ada header Authorization

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}

// TestAlumni_3_GetAll_Endpoint menguji GET /alumni
func TestAlumni_3_GetAll_Endpoint(t *testing.T) {
	// Skenario 1: Berhasil Mendapatkan Semua Alumni (dengan search)
	// <-- TAMBAHAN: Ubah query search agar mencari kedua NIM
	t.Run("Positive - Get All With Search", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/alumni?search=TEST-NIM", nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		meta, ok := respBody["meta"].(map[string]interface{})
		assert.True(t, ok)
		// <-- TAMBAHAN: Harusnya total 2 (1 dari seed, 1 dari create)
		assert.Equal(t, float64(2), meta["total"], "Total data seharusnya 2") 
		data, ok := respBody["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 2)
	})

	// Skenario 2: Berhasil Mendapatkan Semua Alumni (dengan pagination)
	t.Run("Positive - Get All With Pagination", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/alumni?limit=1&page=1", nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		data, ok := respBody["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 1, "Data per halaman seharusnya 1")
	})

	// Skenario 3: Gagal karena Tidak Ada Token (Unauthorized)
	t.Run("Negative - Unauthorized (No Token)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/alumni", nil)
		// Tidak ada header Authorization

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}

// TestAlumni_4_Update_Endpoint menguji PUT /alumni/:id
func TestAlumni_4_Update_Endpoint(t *testing.T) {
	assert.NotEmpty(t, testCreatedAlumniID, "Test Update gagal: testCreatedAlumniID kosong")

	// Skenario 1: Berhasil Update Alumni
	t.Run("Positive - Update Alumni Successfully", func(t *testing.T) {
		reqBody := model.UpdateAlumniRequest{
			Nama:       "Nama Tester Updated", // Data diubah
			Jurusan:    "Sistem Informasi",   // Data diubah
			Angkatan:   2020,
			TahunLulus: 2024,
			Email:      "tester-updated@alumni.com",
			NoTelepon:  "08987654321",
			Alamat:     "Jalan Uji Coba Updated",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PUT", "/api/mg/alumni/"+testCreatedAlumniID, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		data, ok := respBody["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "Nama Tester Updated", data["nama"])
		assert.Equal(t, "Sistem Informasi", data["jurusan"])
	})

	// Skenario 2: Gagal karena ID Tidak Ditemukan (Not Found)
	t.Run("Negative - Not Found ID", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		reqBody := model.UpdateAlumniRequest{Nama: "Update Gagal"}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("PUT", "/api/mg/alumni/"+nonExistentID, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	// Skenario 3: Gagal karena Format ID Salah
	t.Run("Negative - Invalid ID Format", func(t *testing.T) {
		// KIRIM JSON KOSONG AGAR LOLOS BodyParser
		req := httptest.NewRequest("PUT", "/api/mg/alumni/invalid-id", bytes.NewBufferString(`{}`))
		req.Header.Set("Content-Type", "application/json") // <-- Tambahkan content-type
		req.Header.Set("Authorization", "Bearer "+testAuthToken)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode) // Sesuai implementasi repo
	})

	// Skenario 4: Gagal karena Body Request Tidak Valid
	t.Run("Negative - Invalid Request Body", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/mg/alumni/"+testCreatedAlumniID, bytes.NewBufferString(`{"nama":`)) // JSON rusak
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	// Skenario 5: Gagal karena Tidak Ada Token (Unauthorized)
	t.Run("Negative - Unauthorized (No Token)", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/mg/alumni/"+testCreatedAlumniID, nil)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}

// TestAlumni_5_SoftDelete_Endpoint menguji DELETE /alumni/soft-delete/:id
func TestAlumni_5_SoftDelete_Endpoint(t *testing.T) {
	assert.NotEmpty(t, testCreatedAlumniID, "Test SoftDelete gagal: testCreatedAlumniID kosong")

	// Skenario 1: Berhasil Soft Delete
	t.Run("Positive - Soft Delete Successfully", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/mg/alumni/soft-delete/"+testCreatedAlumniID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)

		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, "Alumni berhasil di soft delete", respBody["message"])
	})

	// Skenario 2: Verifikasi Get By ID Gagal Setelah di Soft-Delete (Not Found)
	t.Run("Negative - Get After Soft Delete", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/mg/alumni/"+testCreatedAlumniID+"/", nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode, "Seharusnya 404 Not Found setelah soft delete")
	})

	// Skenario 3: Gagal Menghapus Data yang Sudah Dihapus (Not Found)
	t.Run("Negative - Delete Already Deleted", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/mg/alumni/soft-delete/"+testCreatedAlumniID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode, "Seharusnya 404 Not Found saat menghapus data yg sudah di-soft-delete")
	})

	// Skenario 4: Gagal karena ID Tidak Ditemukan (Not Found)
	t.Run("Negative - Not Found ID", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		req := httptest.NewRequest("DELETE", "/api/mg/alumni/soft-delete/"+nonExistentID, nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	// Skenario 5: Gagal karena Format ID Salah
	t.Run("Negative - Invalid ID Format", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/mg/alumni/soft-delete/invalid-id", nil)
		req.Header.Set("Authorization", "Bearer "+testAuthToken)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode) // Sesuai implementasi repo
	})

	// Skenario 6: Gagal karena Tidak Ada Token (Unauthorized)
	t.Run("Negative - Unauthorized (No Token)", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/mg/alumni/soft-delete/"+testCreatedAlumniID, nil)
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}