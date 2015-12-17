package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/Kutabe/vk"
	"github.com/antonholmquist/jason"
	"github.com/kennygrant/sanitize"
)

func getMusicAsync(wg *sync.WaitGroup, responseURL string, responseArtist string, responseTitle string) {
	fmt.Printf("URL GET:%s\n", responseURL)
	defer wg.Done()
	response, err := http.Get(responseURL)
	if err != nil {
		log.Println(err)
		return
	}
	defer response.Body.Close()
	audio, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}
	filename := sanitize.BaseName(responseArtist+" - "+responseTitle) + ".mp3"
	err = ioutil.WriteFile(filename, audio, 0777)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("FILE DONE:%s\n", filename)
}
func main() {
	var wg sync.WaitGroup
	var login, password string
	maxConnectionCount := flag.Int("maxcon", 10, "Maximum number of download threads")
	flag.Parse()
	fmt.Print("Automatic VK.com music downloader\n")
	fmt.Print("Login(Email or Phone): ")
	fmt.Scanf("%s\n", &login)
	fmt.Print("Password: ")
	fmt.Scanf("%s\n", &password)
	user, err := vk.Auth(login, password)
	if err != nil {
		log.Fatal(err)
	}
	responseJSON, _ := vk.Request("audio.get", nil, user)
	responseObject, _ := jason.NewObjectFromBytes(responseJSON)
	responseArray, _ := responseObject.GetObjectArray("response")
	conCounter := 0
	for _, responseElement := range responseArray {
		responseURL, _ := responseElement.GetString("url")
		responseArtist, _ := responseElement.GetString("artist")
		responseTitle, _ := responseElement.GetString("title")
		conCounter++
		wg.Add(1)
		go getMusicAsync(&wg, responseURL, responseArtist, responseTitle)
		if conCounter > *maxConnectionCount {
			wg.Wait()
			conCounter = 0
		}
	}
	wg.Wait()
}
