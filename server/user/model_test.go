package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&UserModel{})
	assert.NoError(t, err)

	return db
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)

	u := User{
		Username:   ptr("johndoe"),
		FirstName:  ptr("John"),
		LastName:   ptr("Doe"),
		Email:      ptr("john@example.com"),
		Password:   ptr("securepassword"),
		UserStatus: ptr(int32(1)),
	}

	model, err := Create(db, u)
	assert.NoError(t, err)
	assert.Equal(t, "johndoe", model.Username)
}

func TestFindUserByUsername(t *testing.T) {
	db := setupTestDB(t)

	u := User{Username: ptr("janedoe"), Password: ptr("pwd")}
	_, _ = Create(db, u)

	found, err := FindUserByUsername(db, "janedoe")
	assert.NoError(t, err)
	assert.Equal(t, "janedoe", found.Username)
}

func TestUpdateUser(t *testing.T) {
	db := setupTestDB(t)

	_, _ = Create(db, User{Username: ptr("changeme"), Email: ptr("old@e.com")})
	err := Update(db, "changeme", User{Email: ptr("new@e.com")})
	assert.NoError(t, err)

	updated, _ := FindUserByUsername(db, "changeme")
	assert.Equal(t, "new@e.com", updated.Email)
}

func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)

	_, _ = Create(db, User{Username: ptr("tobedeleted")})
	err := Delete(db, "tobedeleted")
	assert.NoError(t, err)

	_, err = FindUserByUsername(db, "tobedeleted")
	assert.Error(t, err)
}

func TestListUsers(t *testing.T) {
	db := setupTestDB(t)

	_, _ = Create(db, User{Username: ptr("a")})
	_, _ = Create(db, User{Username: ptr("b")})

	users, err := ListUsers(db)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func ptr[T any](v T) *T {
	return &v
}
