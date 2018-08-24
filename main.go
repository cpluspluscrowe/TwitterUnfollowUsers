package main

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"net/http"
	"os"
	"time"
)

func doEvery(d time.Duration, f func()) {
	for _ = range time.Tick(d) {
		f()
	}
}

func main() {
	doEvery(300000*time.Millisecond, UnfollowAUserNotFollowingMe)
}

func UnfollowAUserNotFollowingMe() {
	client := getTwitterClient()
	friends, _, err := client.Friends.List(&twitter.FriendListParams{})
	users := friends.Users
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		relationship, _, err := client.Friendships.Show(&twitter.FriendshipShowParams{TargetID: user.ID})
		if err != nil {
			panic(err)
		}
		if !relationship.Target.Following {
			fmt.Println("Unfollowing: ", relationship.Target)
			client.Friendships.Destroy(&twitter.FriendshipDestroyParams{UserID: user.ID})
			break
		}
	}

}

func getClient() *http.Client {
	consumerKey := os.Getenv("consumerKey")
	consumerSecret := os.Getenv("consumerSecret")
	accessToken := os.Getenv("accessToken")
	accessSecret := os.Getenv("accessSecret")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return httpClient
}

func getTwitterClient() *twitter.Client {
	httpClient := getClient()
	client := twitter.NewClient(httpClient)
	return client
}
