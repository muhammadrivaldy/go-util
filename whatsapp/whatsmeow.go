package whatsapp

import (
	"context"

	"github.com/golang/protobuf/proto"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	wProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
)

type whatsMeow struct {
	client *whatsmeow.Client
	device *store.Device
}

func clientWhatsMeow() IClient {

	container, err := sqlstore.New("sqlite3", "file:device.db?_foreign_keys=on", nil)
	if err != nil {
		panic(err)
	}

	device, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	return &whatsMeow{
		client: whatsmeow.NewClient(device, nil),
		device: device,
	}

}

func (w *whatsMeow) Connect(f func(s string)) (IClient, error) {

	if w.client.IsConnected() == false {

		// if the store id is nil, we need to connect it to the smartphone via QR code
		if w.client.Store.ID == nil {

			qrChan, _ := w.client.GetQRChannel(context.Background())

			if err := w.client.Connect(); err != nil {
				return nil, err
			}

			for evt := range qrChan {
				if evt.Event == "code" {
					f(evt.Code)
				}
			}

		} else {

			if err := w.client.Connect(); err != nil {
				return nil, err
			}

		}

	}

	return w, nil

}

func (w *whatsMeow) Closed() {
	w.client.Disconnect()
}

func (w *whatsMeow) Send(phone string, message string) (err error) {

	jid := types.NewJID(phone, types.DefaultUserServer)

	if _, err = w.client.SendMessage(context.Background(), jid, &wProto.Message{
		Conversation: proto.String(message),
	}); err != nil {
		return err
	}

	return

}
