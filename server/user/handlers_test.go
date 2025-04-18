package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// 使用 sqlite :memory: 作为测试数据库
func SetDB(db *gorm.DB) {
	DB = db
}

func setupRouterWithDB(t *testing.T) *chi.Mux {
	db := setupTestDB(t)
	SetDB(db) // 你需要确保你的 handler 能使用这个 DB 实例，或者你可以在 handler 内部注入

	r := chi.NewRouter()
	RegisterRoutes(r)
	return r
}

func TestCreateUserHandler(t *testing.T) {
	r := setupRouterWithDB(t)

	body := User{
		Username:   ptr("apitest"),
		FirstName:  ptr("API"),
		LastName:   ptr("Tester"),
		Email:      ptr("api@example.com"),
		Password:   ptr("123456"),
		UserStatus: ptr(int32(1)),
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetUserByNameHandler(t *testing.T) {
	r := setupRouterWithDB(t)

	// 先创建用户
	createUserForTest(r, t, "getme")

	req := httptest.NewRequest(http.MethodGet, "/user/getme", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "getme")
}

func TestUpdateUserHandler(t *testing.T) {
	r := setupRouterWithDB(t)

	createUserForTest(r, t, "updateme")

	body := User{FirstName: ptr("Updated")}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/user/updateme", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestDeleteUserHandler(t *testing.T) {
	r := setupRouterWithDB(t)

	createUserForTest(r, t, "deleteme")

	req := httptest.NewRequest(http.MethodDelete, "/user/deleteme", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func createUserForTest(r http.Handler, t *testing.T, username string) {
	body := User{
		Username:   ptr(username),
		FirstName:  ptr("Test"),
		Password:   ptr("pass"),
		UserStatus: ptr(int32(1)),
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("createUserForTest failed: %s", resp.Body.String())
	}
}
