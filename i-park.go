package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var lastLogin time.Time = time.Time{}
var currentLoginToken string = ""

func GetAccessToken(device string) (string, error) {
	currentLogin := time.Now().UTC()
	if currentLogin.Sub(lastLogin).Minutes() < 50 {
		return currentLoginToken, nil
	}
	lastLogin = currentLogin
	req, err := http.NewRequest(http.MethodPost, "https://center.hdc-smart.com/v3/auth/login", bytes.NewBufferString("V2"))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Host", "center.hdc-smart.com")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "ionic://localhost")
	req.Header.Set("Authorization", device)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 15_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Status)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("!!!")
	}
	data := make(map[string]string)
	json.Unmarshal([]byte(b), &data)
	currentLoginToken = data["access-token"]
	if currentLoginToken == "" {
		return "", fmt.Errorf("login fail!")
	}
	return currentLoginToken, nil
}

// 차차방 : 3
// 잔누나방 : 4
func SetLight(device string, light bool, room string) error {
	text := "on"
	if light == false {
		text = "off"
	}
	p := "{\"unit\":\"switch1\",\"state\":\"" + text + "\"}"

	req, err := http.NewRequest(http.MethodPut,
		"https://ism.hdc-smart.com/v2/api/features/light/"+string(room)+"/apply",
		bytes.NewBufferString(p))

	if err != nil {
		return err
	}
	token, err := GetAccessToken(device)
	if err != nil {
		return err
	}

	req.Header.Set("Host", "ism.hdc-smart.com")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Origin", "ionic://localhost")
	req.Header.Set("access-token", token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 15_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	respBody, err := io.ReadAll(resp.Body)
	log.Println("light : ", string(respBody))
	if err != nil {
		return err
	}
	return nil
}
