package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

var srv *drive.Service

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	var tok *oauth2.Token
	tokFile := "token.json"
	if _, err := os.Stat(tokFile); os.IsExist(err) {
		tok, err := tokenFromFile(tokFile)
		if err != nil {
			tok = getTokenFromWeb(config)
			saveToken(tokFile, tok)
		}
	} else {
		tok := &oauth2.Token{}
		t := os.Getenv("TOKEN")
		if t == "" {
			log.Println("Ошибка загрузки конфига")
			return nil
		}
		err = json.NewDecoder(strings.NewReader(t)).Decode(tok)
	}

	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", strings.ReplaceAll(authURL, ".metadata.readonly", ""))

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// InitDriveAPI инициализурует апи
func InitDriveAPI() {
	var (
		b   []byte
		err error
	)
	if _, err := os.Stat("credentials.json"); os.IsExist(err) {
		b, err = ioutil.ReadFile("credentials.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}
	} else {
		c := os.Getenv("CREDENTIALS")
		if c == "" {
			log.Println("Не удаеться получить данные доступа")
		}
		b = []byte(c)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err = drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
}

// UploadPhoto загружает фотку на диск
func UploadPhoto(image []byte) string {
	r, err := srv.Files.List().Q("name = 'FotoControll'").
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	if len(r.Files) == 0 {
		fmt.Println("No files found.")
		return ""
	}
	// TODO: имя для фотографии
	f, err := createFile(srv, "", "image/*", bytes.NewBuffer(image), r.Files[0].Id)
	if err != nil {
		fmt.Printf("Error upload file %s", err.Error())
	}
	return f.Id
}

// GetImage получение картинки по его id
func GetImage(idImage string) []byte {
	r, err := srv.Files.Get(idImage).Download()
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil
	}
	defer r.Body.Close()
	image, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil
	}
	return image
}

func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := service.Files.Create(f).Media(content).Do()

	if err != nil {
		log.Println("Could not create file: " + err.Error())
		return nil, err
	}

	return file, nil
}
