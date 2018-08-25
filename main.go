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
	UnfollowAUserNotFollowingMe()
	doEvery(86400*time.Second, UnfollowAUserNotFollowingMe)
}

func GetNextFriendsCursor(client *twitter.Client, cursor int64) int64 {
	var friendParams *twitter.FriendListParams
	if cursor > -1 {
		friendParams = &twitter.FriendListParams{}
	} else {
		friendParams = &twitter.FriendListParams{Cursor: cursor}
	}
	friends, _, err := client.Friends.List(friendParams)
	if err != nil {
		panic(err)
	}
	friendsData := friends
	nextCursor := friendsData.NextCursor
	return nextCursor
}

func GetOlderFriends(client *twitter.Client) []twitter.User {
	var cursor int64 = -1
	for i := 0; i < 4; i++ {
		cursor = GetNextFriendsCursor(client, cursor)
		time.Sleep(10 * time.Millisecond)
		fmt.Println(cursor)
	}
	friendParams := &twitter.FriendListParams{Cursor: cursor}
	friends, _, err := client.Friends.List(friendParams)
	if err != nil {
		panic(err)
	}
	users := friends.Users
	return users
}

func UnfollowAUserNotFollowingMe() {
	client := getTwitterClient()
	users := GetOlderFriends(client)
	for _, user := range users {
		relationship, _, err := client.Friendships.Show(&twitter.FriendshipShowParams{TargetID: user.ID})
		if err != nil {
			panic(err)
		}
		if !relationship.Target.Following {
			fmt.Println("Unfollowing: ", relationship.Target)
			client.Friendships.Destroy(&twitter.FriendshipDestroyParams{UserID: user.ID})
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Finished unfollowing method")
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
