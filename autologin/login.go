//Function name always be Caps
package autologin

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/stealth"
)

type Cookie []struct {
	Path     string `json:"path"`
	Domain   string `json:"domain"`
	Value    string `json:"value"`
	Name     string `json:"name"`
	HttpOnly bool   `json:"httpOnly"`
}

func Login(username string, password string) {
	page := rod.New().MustConnect()
	stealth_page := stealth.MustPage(page)
	var cookie Cookie

	stealth_page.MustNavigate("https://www.instagram.com/accounts/login")
	time.Sleep(3 * time.Second)

	insert_username := stealth_page.MustElementX("//*[@id=\"loginForm\"]/div/div[1]/div/label/input")
	// Input username
	user_name := username

	for _, rune_ := range user_name {
		insert_username.MustInput(string(rune_))
		time.Sleep(400 * time.Millisecond)

	}

	// Input password
	insert_password := stealth_page.MustElementX("//*[@id=\"loginForm\"]/div/div[2]/div/label/input")
	user_password := password

	for _, rune_ := range user_password {
		insert_password.MustInput(string(rune_))
		time.Sleep(400 * time.Millisecond)

	}
	time.Sleep(2 * time.Second)

	stealth_page.MustElementX("//*[@id=\"loginForm\"]/div/div[3]/button").MustPress(input.Enter)

	time.Sleep(6 * time.Second)

	//Storing cookies to json file

	cookies := stealth_page.Browser().MustGetCookies()
	data, _ := json.Marshal(cookies)

	json.Unmarshal([]byte(data), &cookie)
	newdata, _ := json.MarshalIndent(cookie, "", "")

	file, ferr := os.OpenFile("cookie.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if ferr == nil {
		file.WriteString(string(newdata))
	}
	fmt.Println("***************Login Successful******************")

	stealth_page.Close()
}
