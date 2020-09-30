package util

import (
	"encoding/json"
	"github.com/devops-salt/src/message"
	"io/ioutil"
	"net/http"
	"strings"
)

func Post(url string, body string) (int,error) {
	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func LoadTask(url string) (*message.Task, error){
	body, err := Get(url)
	if err != nil {
		return nil, err
	}
	task := &message.Task{}
	if err := json.Unmarshal(body, task); err != nil {
		return task, err
	}
	return task, nil
}


func Get(url string) ([]byte, error){
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
