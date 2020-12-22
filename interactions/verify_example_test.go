package interactions_test

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/bsdlp/discord-interactions-go/interactions"
)

func ExampleVerify() {
	hexEncodedDiscordPubkey := "a43f74054d052d43c3ed90e07c8e3270690826d1fd38eda3a42e91ff38c2482b"
	discordPubkey, err := hex.DecodeString(hexEncodedDiscordPubkey)
	if err != nil {
		// handle error
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		verified := interactions.Verify(r, ed25519.PublicKey(discordPubkey))
		if !verified {
			http.Error(w, "signature mismatch", http.StatusUnauthorized)
			return
		}

		defer r.Body.Close()
		var data interactions.Data
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			// handle error
		}

		// respond to ping
		if data.Type == interactions.Ping {
			_, err := w.Write([]byte(`{"type":1}`))
			if err != nil {
				// handle error
			}
			return
		}

		// handle command
		response := &interactions.InteractionResponse{
			Type: interactions.ChannelMessage,
			Data: &interactions.InteractionApplicationCommandCallbackData{
				Content: "got your message kid",
			},
		}

		var responsePayload bytes.Buffer
		err = json.NewEncoder(&responsePayload).Encode(response)
		if err != nil {
			// handle error
		}

		_, err = http.Post(data.ResponseURL(), "application/json", &responsePayload)
		if err != nil {
			// handle err
		}
	})
}
