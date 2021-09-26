package Toyoutube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const (
	CredentialFilePath = "../../lib/client_secret_530945973561-lbmkdisri0iat9chvrnt5ehdocupk9rm.apps.googleusercontent.com.json"
	TokenFilePath      = "../../lib/token.json"
	MaxResult          = 50

	OldUrl = "https://old.url.com"
	NewUrl = "https://new.url.com"
)

// token をローカルもしくは新規にWebから取得
func getToken(config *oauth2.Config) *oauth2.Token {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(TokenFilePath)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(TokenFilePath, tok)
	}
	return tok
}

// Webにアクセスし OAuth 2 認証をおこない token を取得
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("下記のリンクへ飛び認証コードを発行し、その文字列をここに貼り付けて Enter を押してください: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// token をローカルファイルから取得
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

// token をローカルファイルとして保存
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// service を取得
func getService() *youtube.Service {
	b, err := ioutil.ReadFile(CredentialFilePath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, youtube.YoutubeScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	tok := getToken(config)

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, tok)))
	if err != nil {
		log.Fatalf("Unable to create people Client %v", err)
	}

	return service
}

func Update(URL string) {
	service := getService()
	counter := 0

	//searchCall := service.Search.List([]string{"id"}).ForMine(true).Type("video").MaxResults(1)
	//searchResponse, err := searchCall.Do()

	videoId := URL[31:]

	videoListCall := service.Videos.List([]string{"id", "snippet"}).Id(videoId).MaxResults(1)

	videoListResponse, err := videoListCall.Do()
	if err != nil {
		log.Fatalf("Error making API call to list videos: %v", err.Error())
	}
	title := videoListResponse.Items[0].Snippet.Title
	description := videoListResponse.Items[0].Snippet.Description
	if strings.Contains(description, OldUrl) {
		fmt.Printf("%v (%v)\n", title, videoId)
		videoListResponse.Items[0].Snippet.Description = strings.Replace(description, OldUrl, NewUrl, -1)
		service.Videos.Update([]string{"snippet"}, videoListResponse.Items[0]).Do()
	}

	fmt.Println(counter)
}
