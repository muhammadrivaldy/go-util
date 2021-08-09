package util

import (
	"context"
	"fmt"
	"runtime"
	"testing"
)

var tTeleService TeleService

func init() {
	var err error
	tTeleService, err = NewTele("1436691401:AAFbIzoj6Ymtc8EYrHhLEIvqmF7a9FsRUyQ", -352711666)
	if err != nil {
		fmt.Println("Error", err.Error())
	}
}

// TestSendError how to run this process
// go test -v -run=TestSendError
func TestSendError(t *testing.T) {
	_, path, line, _ := runtime.Caller(0)
	err := tTeleService.SendError(context.Background(), path, line, "Bagas yang punya kawasan Jaktim dan sekitarnya")
	if err != nil {
		t.Error(err)
		return
	}
}
