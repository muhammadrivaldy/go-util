package whatsapp_test

import (
	"os"
	"testing"

	"github.com/mdp/qrterminal/v3"
	"github.com/muhammadrivaldy/go-util/whatsapp"
)

func TestSend(t *testing.T) {

	client, err := whatsapp.NewClient().Connect(func(s string) { qrterminal.GenerateHalfBlock(s, qrterminal.L, os.Stdout) })
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 100; i++ {
		if err := client.Send("6287723137610", "Just a testing!"); err != nil {
			t.Error(err)
		}
	}

}
