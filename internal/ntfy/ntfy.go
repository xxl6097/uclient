package ntfy

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type NtfyEventData struct {
	Id    string `json:"id,omitempty"`
	Time  int    `json:"time,omitempty"`
	Event string `json:"event,omitempty"`
	Topic string `json:"topic,omitempty"`
}
type Ntfy struct {
}

func subscribe(address, topic string, username, password string, fn func(string)) error {
	//req, e := http.NewRequest("GET", fmt.Sprintf("%s/%s/json?poll=1", address, topic), nil)
	req, e := http.NewRequest("GET", fmt.Sprintf("%s/%s/json", address, topic), nil)
	if e != nil {
		return e
	}
	req.SetBasicAuth(username, password)
	//req.Header.Set("Authorization", "Basic tk_trk4agho2")
	req.Header.Set("UUXIA", "Basic tk_trk4agho2")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	scanner := bufio.NewScanner(resp.Body)
	if e1 := scanner.Err(); e1 != nil {
		return e1
	}
	for scanner.Scan() {
		if fn != nil {
			fn(scanner.Text())
		}
	}
	return scanner.Err()
}

func Subscribe(address, topic string, username, password string, fn func(string)) {
	err := subscribe(address, topic, username, password, fn)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Second * 10)
		Subscribe(address, topic, username, password, fn)
	}
}
