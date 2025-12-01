package tests

import (
	"online-learning-platform/pkg/utils"
	"testing"
)

func TestPassword(t *testing.T) {
	password := "teacher123"
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hashedPassword)
}
