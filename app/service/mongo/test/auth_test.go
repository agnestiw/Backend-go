package test

import (
	"bytes"
	"encoding/json"
	"latihan2/app/model"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestLogin_Endpoint(t *testing.T) {

	// --- Skenario 1: Login Berhasil ---
	t.Run("Login Berhasil", func(t *testing.T) {
		reqBody := model.LoginRequest{
			Username: testSeededUsername,
			Password: testSeededPassword,
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/mg/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := testApp.Test(req, -1) 
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode, "Status code seharusnya 200 OK")
		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		data, ok := respBody["data"].(map[string]interface{})
		assert.True(t, ok, "Key 'data' seharusnya ada")
		assert.NotEmpty(t, data["token"], "Token seharusnya tidak kosong")
		userResp, ok := data["user"].(map[string]interface{})
		assert.True(t, ok, "Key 'user' seharusnya ada")
		assert.Equal(t, testSeededUsername, userResp["username"], "Username di respons tidak sesuai")
	})

	// --- Skenario 2: Password Salah ---
	t.Run("Password Salah", func(t *testing.T) {
		reqBody := model.LoginRequest{
			Username: testSeededUsername,
			Password: "password-salah-pasti", 
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/mg/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode, "Status code seharusnya 401 Unauthorized")
		var respBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Contains(t, respBody["error"], "Username atau password salah", "Pesan error tidak sesuai")
	})

	// --- Skenario 3: User Tidak Ditemukan ---
	t.Run("User Tidak Ditemukan", func(t *testing.T) {
		reqBody := model.LoginRequest{
			Username: "user_pasti_tidak_ada", 
			Password: testSeededPassword,
		}
		bodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/api/mg/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode, "Status code seharusnya 401 Unauthorized")
	})

	// --- Skenario 4: Body Request Tidak Valid ---
	t.Run("Body Request Tidak Valid", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/mg/login", bytes.NewBufferString(`{"username": "admin",`)) // JSON rusak
		req.Header.Set("Content-Type", "application/json")
		resp, err := testApp.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "Status code seharusnya 400 Bad Request")
	})
}