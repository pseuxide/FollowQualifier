package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/schollz/progressbar"
)

var (
	api               *anaconda.TwitterApi
	ConsumerKey             = ""
	ConsumerSecret          = ""
	AccessToken             = ""
	AccessTokenSecret       = ""
	listName                = ""
	listId            int64 = 123
)

func main() {
	api = anaconda.NewTwitterApiWithCredentials(AccessToken, AccessTokenSecret, ConsumerKey, ConsumerSecret)
	removeTarget := getRemoveTarget()
	fmt.Printf("[+]%v targets leave\n", len(removeTarget))
	remove(removeTarget)
}

//remove unfollows accounts passed by parent
func remove(t []int64) {
	bar := progressbar.New(1000)
	for r, i := range t {
		bar.Add(1)
		_, err := api.UnfollowUserId(i)
		//		fmt.Printf("[-]%v people remove->%v\n", r+1, u.Name)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(1000 * time.Millisecond)
		if r == 999 { //if number of removal hits 1000, BREAK
			break
		}
	}
}

//getPreserveTargets returns allay of accounts which isn't supposed to be removed according to the particular list in Twitter
func getPreserveTargets(listName string, listId int64) []int64 {
	var userIds = []int64{}
	var NextCursor int64
	var flag bool
	for NextCursor != -1 {
		if flag == false {
			NextCursor = -1
			flag = true
		}
		v := url.Values{"cursor": {strconv.FormatInt(NextCursor, 10)}}
		userCursor, err := api.GetListMembers(listName, listId, v)
		if err != nil {
			fmt.Println(err)
		}
		for _, user := range userCursor.Users {
			userIds = append(userIds, user.Id)
		}
		NextCursor = userCursor.Next_cursor
	}
	return userIds
}

//getFollowing retrieves all your following user's ids and returns them
func getFollowing() []int64 {
	nextCursor := "-1"
	friends := []int64{}
	for i := 0; i < 15; i++ {
		v := url.Values{"cursor": {nextCursor}}
		f, err := api.GetFriendsIds(v)
		if err != nil {
			log.Println(err)
		}
		friends = append(friends, f.Ids...)
		nextCursor = f.Next_cursor_str
		if nextCursor == "0" {
			break
		}
	}
	return friends
}

//getRemoveTarget reveals who to remove clearly, and returns them
func getRemoveTarget() []int64 {
	followers := getFollowing()
	preserveTargets := getPreserveTargets(listName, listId)

	filter := map[int64]int{}
	for _, f := range followers {
		filter[f] = 1
	}
	for _, pt := range preserveTargets {
		_, exist := filter[pt]
		if exist {
			filter[pt] = 2
		}
	}
	var result = []int64{}
	for f, i := range filter {
		if i == 1 {
			result = append(result, f)
		}
	}
	return result
}

/*
func getRemoveTarget() []int64 {
	nextCursor := "-1"
	friends := []int64{}
	for i := 0; i < 15; i++ {
		v := url.Values{"cursor": {nextCursor}}
		f, err := api.GetFriendsIds(v)
		if err != nil {
			log.Println(err)
		}
		friends = append(friends, f.Ids...)
		nextCursor = f.Next_cursor_str
		if nextCursor == "0" {
			break
		}
	}
	nextCursor = "-1"
	followers := []int64{}
	for i := 0; i < 15; i++ {
		v := url.Values{"cursor": {nextCursor}}
		fed, err := api.GetFollowersIds(v)
		if err != nil {
			log.Println(err)
		}
		followers = append(followers, fed.Ids...)
		nextCursor = fed.Next_cursor_str
		if nextCursor == "0" {
			break
		}
	}
	return filter(friends, followers)
}

*/
//lhsの中から、rhsと共通しないものを吐き出す関数l
//lhsがフォロー・rhsがフォロワー

//example:
//lhs := []int64{1, 3, 5, 7, 9}
//rhs := []int64{3, 5}
//filter(lhs, rhs) -> 1, 7, 9
/*
func filter(lhs, rhs []int64) []int64 {
	m := map[int64]int{}

	for _, v := range lhs {
		if _, ok := m[v]; !ok {
			m[v] = 1
		}
	}

	for _, v := range rhs {
		if _, ok := m[v]; ok {
			m[v] = 2
		}
	}

	var ret []int64

	for i := range m {
		if m[i] == 1 {
			ret = append(ret, i)
		}
	}
	return ret
}
*/
