package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	autologin "InstaOsint/autologin"
	fetch "InstaOsint/getlogindata"

	"github.com/briandowns/spinner"
)

// Read the cookies which are used to login
func cookieReader() string {
	f, err := os.ReadFile("./cookie.json")
	var getLoginCookies string

	if err != nil {
		fmt.Println(err.Error())
	} else {
		var cookie []map[string]interface{}
		// fmt.Println(string(f))

		json.Unmarshal([]byte(f), &cookie)

		for _, v := range cookie {
			if v["name"] == "sessionid" {
				getLoginCookies = v["name"].(string) + "=" + v["value"].(string)
			}
		}
	}
	return getLoginCookies
}

func resp(w http.ResponseWriter, r *http.Request) {

	var uri url.Values
	var conf map[string]interface{}
	uri = r.URL.Query()

	args := strings.Split(uri.Get("q"), ":")

	// Reading Configuration here

	config, _ := os.ReadFile("./config.json") // read username and password
	json.Unmarshal([]byte(config), &conf)

	if conf["username"] == "" && conf["password"] == "" {
		fmt.Print("\n")
		fmt.Println("Error while reading the username and password from configuration File !!!! Please enter the configuration correctly then execute program again.")
		fmt.Print("\n")
		os.Exit(1)
	} else {
		intercookie := cookieReader()

		data, err := fetch.Topsearches(intercookie, args[0])
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			if data == nil {
				autologin.Login(conf["username"].(string), conf["password"].(string))
				newcookies := cookieReader()
				data, _ := fetch.Topsearches(newcookies, args[0])
				fmt.Println("**********Login Cookie Expired************")
				fmt.Print("\n")

				d, _ := json.Marshal(data)
				fmt.Print("\n")

				fmt.Println(data)
				w.Write([]byte(d))
			} else {
				fmt.Println("**********Session Cookies************")
				fmt.Print("\n")

				fmt.Println(data)
				d, _ := json.Marshal(data)
				fmt.Print("\n")

				w.Write([]byte(d))
			}
		}
	}
}

func main() {
	fmt.Println("Reading Configuration.....")

	s := spinner.New(spinner.CharSets[36], 100*time.Millisecond) // Build our new spinner
	s.Color("cyan", "bold")

	fmt.Println("Server is starting")
	s.Start() // Start the spinner
	time.Sleep(4 * time.Second)

	http.HandleFunc("/social", resp)
	s.FinalMSG = "Server is listening on localhost:8000 !!\n"
	s.Stop()
	http.ListenAndServe(":8000", nil)

}
