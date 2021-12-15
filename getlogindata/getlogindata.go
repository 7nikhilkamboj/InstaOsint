//Function name always be Caps

package getlogindata

import (
	conversion "InstaOsint/conversion"
	"encoding/json"
	"fmt"

	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

type UserDetails struct {
	Pk             int         `json:"pk"`
	Username       interface{} `json:"username"`
	FollowerCount  interface{} `json:"follower_count"`
	FollowingCount interface{} `json:"following_count"`
	FullName       interface{} `json:"full_name"`
	ProfilePicURL  interface{} `json:"profile_pic_url"`
	Biography      interface{} `json:"biography"`
	MediaCount     interface{} `json:"media_count"`
	IsVerified     interface{} `json:"is_verified"`
}

func Parser(data []byte, intercookie string) []UserDetails {
	/*
		d, _ := json.Marshal(response.Body)
		fmt.Println(response.Body)
	*/
	var parsed map[string]interface{}

	Complete := []UserDetails{}

	json.Unmarshal([]byte(data), &parsed)

	d := parsed["users"]
	var pk []interface{}
	for k := range d.([]interface{}) {
		test := d.([]interface{})[k].(map[string]interface{})["user"].(map[string]interface{})["pk"]
		// fmt.Println(test)
		pk = append(pk, test)
	}
	for i := range pk {
		if i <= 9 {
			// returndata = append(returndata, GetLogin(intercookie, pk[i]))
			Complete = append(Complete, GetLogin(intercookie, pk[i]))
		}
	}
	return Complete
}

func Topsearches(intercookie string, args string) ([]UserDetails, error) {
	//https://www.instagram.com/web/search/topsearch/?context=blended&query={usersearch}

	url_topsearch := "https://www.instagram.com/web/search/topsearch/?context=blended&query=" + url.QueryEscape(args) + ""
	method := "GET"
	var d []UserDetails
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url_topsearch, nil)
	if err != nil {
		fmt.Println("TopSearch Error:", err)

	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:95.0) Gecko/20100101 Firefox/95.0")
	req.Header.Add("Cookie", intercookie)
	req.Header.Add("X-IG-App-ID", "936619743392459")
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("TopSearch Error:", err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	resp_body := string(data)

	match, _ := regexp.MatchString("Login • Instagram", resp_body)
	if match {
		fmt.Println("Session Expired")
		fmt.Println("Loggin again******")
	} else {
		d = Parser(data, intercookie)

	}
	return d, err
}

func GetLogin(intercookie string, pk interface{}) UserDetails {
	// url := "https://www.instagram.com/virat.kohli/?__a=1"
	info := "https://i.instagram.com/api/v1/users/" + pk.(string) + "/info/"
	method := "GET"
	var userdata map[string]interface{}
	// var complData []interface{}
	var userdetail UserDetails

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, info, nil)
	if err != nil {
		fmt.Println("Login Error:", err)

	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:95.0) Gecko/20100101 Firefox/95.0")
	req.Header.Add("Cookie", intercookie)
	req.Header.Add("X-IG-App-ID", "936619743392459")
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Login Error:", err)
	}
	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	json.Unmarshal([]byte(data), &userdata)

	d := userdata["user"]

	da := d.(map[string]interface{})["pk"].(float64) // Conversion of userid to roundInt

	// complData = append(complData, conversion.RoundInt(da), d.(map[string]interface{})["username"], conversion.NearestThousandFormat(d.(map[string]interface{})["follower_count"].(float64)), conversion.NearestThousandFormat(d.(map[string]interface{})["following_count"].(float64)), d.(map[string]interface{})["full_name"], d.(map[string]interface{})["profile_pic_url"], d.(map[string]interface{})["biography"], d.(map[string]interface{})["media_count"], d.(map[string]interface{})["is_verified"])

	userdetail.Pk = conversion.RoundInt(da)
	userdetail.Username = d.(map[string]interface{})["username"]
	userdetail.FollowerCount = conversion.NearestThousandFormat(d.(map[string]interface{})["follower_count"].(float64))
	userdetail.FollowingCount = conversion.NearestThousandFormat(d.(map[string]interface{})["following_count"].(float64))
	userdetail.FullName = d.(map[string]interface{})["full_name"]
	userdetail.ProfilePicURL = d.(map[string]interface{})["profile_pic_url"]
	userdetail.Biography = d.(map[string]interface{})["biography"]
	userdetail.MediaCount = d.(map[string]interface{})["media_count"]
	userdetail.IsVerified = d.(map[string]interface{})["is_verified"]

	return userdetail

}
