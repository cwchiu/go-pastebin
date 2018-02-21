package pastebin

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Pastebin represents an instance of the pastebin service.
type Pastebin struct {
	Key     string
	UserKey string
}

// StripURL returns the paste ID from a pastebin URL.
func (p Pastebin) StripURL(url string) string {
	return strings.Replace(url, "https://pastebin.com/", "", -1)
}

// WrapID returns the pastebin URL from a paste ID.
func (p Pastebin) WrapID(id string) string {
	return "https://pastebin.com/" + id
}

func (p *Pastebin) Login(username, password string) error {
	data := url.Values{}
	// Required values.
	data.Set("api_dev_key", p.Key)
	data.Set("api_user_name", username)
	data.Set("api_user_password", password)

	resp, err := http.PostForm("https://pastebin.com/api/api_login.php", data)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("pastebin get failed")
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodyStr := string(respBody)

	if strings.Contains(bodyStr, "Bad API request,") {
		return errors.New(bodyStr)
	}

	p.UserKey = bodyStr

	return nil
}

// 張貼
func (p Pastebin) Put(paste Paste) (id string, err error) {
	data := url.Values{}
	data.Set("api_dev_key", p.Key)
	data.Set("api_option", "paste")

	if paste.MemberOnly == true && p.UserKey != "" {
		data.Set("api_user_key", p.UserKey)
	}

	data.Set("api_paste_code", paste.Code)
	data.Set("api_paste_name", paste.Title)

	if paste.Private == "" {
		paste.Private = "2"
	}
	data.Set("api_paste_private", paste.Private)

	if paste.ExpireDate == "" {
		paste.ExpireDate = "N"
	}
	data.Set("api_paste_expire_date", paste.ExpireDate)

	if paste.FormatShort != "" {
		data.Set("api_paste_format", paste.FormatShort)
	}

	resp, err := http.PostForm("https://pastebin.com/api/api_post.php", data)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("pastebin get failed")
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return p.StripURL(string(respBody)), nil
}

// Get returns the text inside the paste identified by ID.
func (p Pastebin) Get(id string) (text string, err error) {
	resp, err := http.Get("https://pastebin.com/raw.php?i=" + id)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("pastebin get failed")
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}

func (p Pastebin) ListByUser(limit int) (*PasteList, error) {
	body, err := p.UserRequest("https://pastebin.com/api/api_post.php", map[string]string{
		"api_results_limit": strconv.Itoa(limit),
		"api_option":        "list",
	})
	if err != nil {
		return nil, err
	}

	var ret PasteList
	err = xml.Unmarshal([]byte(fmt.Sprintf("<data>%s</data>", body)), &ret)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

// 取得用戶貼文
func (p Pastebin) GetByUser(id string) (text string, err error) {
	body, err := p.UserRequest("https://pastebin.com/api/api_raw.php", map[string]string{
		"api_paste_key": id,
		"api_option":    "show_paste",
	})
	if err != nil {
		return "", err
	}

	return body, nil
}

// 刪除用戶的貼文
func (p Pastebin) DelByUser(id string) (err error) {
	body, err := p.UserRequest("https://pastebin.com/api/api_post.php", map[string]string{
		"api_paste_key": id,
		"api_option":    "delete",
	})
	if err != nil {
		return err
	}

	if strings.Contains(body, "Paste Removed") {
		return nil
	}
	return errors.New(body)
}

// 用戶資訊
func (p Pastebin) InfoByUser() (*User, error) {
	body, err := p.UserRequest("https://pastebin.com/api/api_post.php", map[string]string{
		"api_option": "userdetails",
	})
	if err != nil {
		return nil, err
	}

	user := User{}
	err = xml.Unmarshal([]byte(body), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p Pastebin) UserRequest(api_gateway string, options map[string]string) (text string, err error) {
	data := url.Values{}
	data.Set("api_dev_key", p.Key)
	data.Set("api_user_key", p.UserKey)
	for k, v := range options {
		data.Set(k, v)
	}

	resp, err := http.PostForm(api_gateway, data)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("User Request failed")
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	bodyStr := string(respBody)
	if strings.Contains(bodyStr, "Bad API request,") {
		return "", errors.New(bodyStr)
	}

	return bodyStr, nil
}
