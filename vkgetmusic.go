package main
import (
    "fmt"
    "io/ioutil"
    "log"
    "sync"
    "net/http"
    "github.com/Kutabe/vk"
    "github.com/antonholmquist/jason"
    "github.com/kennygrant/sanitize"
)
func getMusicAsync(wg *sync.WaitGroup,responseURL string, responseArtist string, responseTitle string) {
    defer wg.Done()
     response, err := http.Get(responseURL)
     if err != nil  {
         log.Print(err)
         return
     }
     defer response.Body.Close()
     audio, _ := ioutil.ReadAll(response.Body)
     filename := sanitize.BaseName(responseArtist+" - "+responseTitle)+".mp3"
     ioutil.WriteFile(filename, audio, 0777)
     fmt.Printf("URL GET:%s to file:%s\n",responseURL,filename)
}
func main() {
    var wg sync.WaitGroup
    var login, password string
    fmt.Print("Login: ")
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
    for _, responseElement := range responseArray {
        responseURL, _ := responseElement.GetString("url")
        responseArtist, _ := responseElement.GetString("artist")
        responseTitle, _ := responseElement.GetString("title")

        fmt.Printf("URL%s \n", responseURL)
        wg.Add(1)
        go getMusicAsync(&wg,responseURL,responseArtist,responseTitle)
    }
    wg.Wait()
}
