package apns

import (
	"encoding/json"
	"testing"
)

func TestAlertMarshal(t *testing.T) {
	{
		alert := Alert{}
		alert.LockKey = "GAME_PLAY_REQUEST_FORMAT"
		alert.LockArgs = []string{"Jenna", "Frank"}
		j, err := json.Marshal(alert)
		if err != nil {
			t.Fatalf("can't marshal to json: %s", err)
		}
		if got, expect := string(j), `{"loc-key":"GAME_PLAY_REQUEST_FORMAT","loc-args":["Jenna","Frank"]}`; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		alert := Alert{}
		alert.Body = "Bob wants to play poker"
		alert.ActionLockKey = "PLAY"
		j, err := json.Marshal(alert)
		if err != nil {
			t.Fatalf("can't marshal to json: %s", err)
		}
		if got, expect := string(j), `{"body":"Bob wants to play poker","action-loc-key":"PLAY"}`; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}
}

func TestApsMarshal(t *testing.T) {
	{
		aps := Aps{}
		aps.Alert = Alert{
			LockKey:  "GAME_PLAY_REQUEST_FORMAT",
			LockArgs: []string{"Jenna", "Frank"},
		}
		j, err := json.Marshal(aps)
		if err != nil {
			t.Fatalf("can't marshal to json: %s", err)
		}
		if got, expect := string(j), `{"alert":{"loc-key":"GAME_PLAY_REQUEST_FORMAT","loc-args":["Jenna","Frank"]}}`; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		aps := Aps{}
		aps.Alert.Body = "Message received from Bob"
		j, err := json.Marshal(aps)
		if err != nil {
			t.Fatalf("can't marshal to json: %s", err)
		}
		if got, expect := string(j), `{"alert":{"body":"Message received from Bob"}}`; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		aps := Aps{}
		aps.Alert = Alert{
			Body:          "Bob wants to play poker",
			ActionLockKey: "PLAY",
		}
		aps.Badge = 5
		j, err := json.Marshal(aps)
		if err != nil {
			t.Fatalf("can't marshal to json: %s", err)
		}
		if got, expect := string(j), `{"alert":{"body":"Bob wants to play poker","action-loc-key":"PLAY"},"badge":5}`; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		aps := Aps{}
		aps.Alert.Body = "You got your emails."
		aps.Badge = 9
		aps.Sound = "bingbong.aiff"
		j, err := json.Marshal(aps)
		if err != nil {
			t.Fatalf("can't marshal to json: %s", err)
		}
		if got, expect := string(j), `{"alert":{"body":"You got your emails."},"badge":9,"sound":"bingbong.aiff"}`; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		aps := Aps{}
		j, err := json.Marshal(aps)
		if err != nil {
			t.Fatalf("can't marshal to json: %s", err)
		}
		if got, expect := string(j), `{"alert":{}}`; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}
}

func TestPayloadMarshal(t *testing.T) {
	{
		payload := Payload{}
		payload.Aps.Alert.Body = "Message received from Bob"
		payload.SetCustom("acme2", []string{"bang", "whiz"})
		j, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("can't marshal to json: %s", err)
		}
		if got, expect := string(j), `{"acme2":["bang","whiz"],"aps":{"alert":{"body":"Message received from Bob"}}}`; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		payload := Payload{}
		payload.Aps.Alert.Body = "You got your emails."
		payload.Aps.Badge = 9
		payload.Aps.Sound = "bingbong.aiff"
		payload.SetCustom("acme1", "bar")
		payload.SetCustom("acme2", 42)
		j, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("can't marshal to json: %s", err)
		}
		if got, expect := string(j), `{"acme1":"bar","acme2":42,"aps":{"alert":{"body":"You got your emails."},"badge":9,"sound":"bingbong.aiff"}}`; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		payload := Payload{}
		payload.SetCustom("acme2", []int{5, 8})
		j, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("can't marshal to json: %s", err)
		}
		if got, expect := string(j), `{"acme2":[5,8],"aps":{"alert":{}}}`; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		payload := Payload{}
		payload.Aps.Alert = Alert{
			LockKey:  "GAME_PLAY_REQUEST_FORMAT",
			LockArgs: []string{"Jenna", "Frank"},
		}
		payload.Aps.Sound = "chime"
		payload.SetCustom("acme", "foo")
		j, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("can't marshal to json: %s", err)
		}
		if got, expect := string(j), `{"acme":"foo","aps":{"alert":{"loc-key":"GAME_PLAY_REQUEST_FORMAT","loc-args":["Jenna","Frank"]},"sound":"chime"}}`; got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}
}
