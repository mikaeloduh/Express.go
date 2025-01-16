package e2e

import (
	"os"
	"testing"
)

var (
	userService *UserService
)

func setup() {
	userService = NewUserService()
	userService.CreateUser("correctName", "correctEmail@example.com", "correctPassword")
}

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	os.Exit(code)
}
