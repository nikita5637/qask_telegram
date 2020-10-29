package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Question struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Comment  string `json:"comment"`
}

func GetQuestion(user *User, qaskIP, qaskPort string) *Question {
	type request struct {
		TgId int64  `json:"tgId"`
		From string `json:"from"`
	}
	re := &request{}
	re.TgId = user.UserID()
	re.From = "telegram"

	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(re)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%s/questions", qaskIP, qaskPort), b)

	client := http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	var q Question

	err = json.Unmarshal(data, &q)
	if err != nil {
		return nil
	}

	return &q
}
