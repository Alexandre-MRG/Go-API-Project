package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test si connexion réussie au serveur blockchain
// Test si pas de connexion

type newPlayerTest struct {
	playerArg          Player
	expectedStatusCode int
}

var newPlayerTests = []newPlayerTest{
	{
		Player{
			Username: "usertest",
			Password: "passtest",
			Pin_code: "123456",
		},
		200,
	},
	{
		Player{
			Username: "usertest",
			Password: "PassTest57%", // doit accepter les caractères spéciaux dans le mot de passe
			Pin_code: "123456",
		},
		200,
	},
	{
		Player{
			Username: "us#~*lk", // user contient des caractères spéciaux
			Password: "passtest",
			Pin_code: "123456",
		},
		400,
	},
	{
		Player{
			Username: "usertest",
			Password: "pa", // mot de passe trop court (6-32 characters)
			Pin_code: "123456",
		},
		400,
	},
	{
		Player{
			Username: "usertest",
			Password: "passtest",
			Pin_code: "1234", // pin trop court (6 digits (0-9))
		},
		400,
	},
	{
		Player{},
		400,
	},
}

func TestHandleNewPlayer(t *testing.T) {

	for _, test := range newPlayerTests {

		jsonRequest, err := json.Marshal(test.playerArg)
		if err != nil {
			LogError(SERIALIZE_ERROR_MESSAGE, err)
		}
		req := httptest.NewRequest(http.MethodPost, "/players/create", bytes.NewBuffer(jsonRequest))
		w := httptest.NewRecorder()

		HandleNewPlayer(w, req)

		res := w.Result()
		/* 		data, _ := ioutil.ReadAll(res.Body)

		   		newPlayer := Player{}
		   		if err := json.Unmarshal(data, &newPlayer); err != nil {
		   			LogError(DESERIALIZE_ERROR_MESSAGE, err)
		   			w.WriteHeader(http.StatusBadRequest)
		   			return
		   		} */

		if res.StatusCode != test.expectedStatusCode {
			t.Errorf("Output %d not equal to expected %d", res.StatusCode, test.expectedStatusCode)
		}
	}
}

func BenchmarkHandleNewPlayer(b *testing.B) {

	for i := 0; i < b.N; i++ {
		playerBench := Player{
			Username: "usertest",
			Password: "passtest",
			Pin_code: "123456",
		}

		jsonRequest, err := json.Marshal(playerBench)
		if err != nil {
			LogError(SERIALIZE_ERROR_MESSAGE, err)
		}

		req := httptest.NewRequest(http.MethodPost, "/players/create", bytes.NewBuffer(jsonRequest))
		w := httptest.NewRecorder()

		HandleNewPlayer(w, req)
	}
}
