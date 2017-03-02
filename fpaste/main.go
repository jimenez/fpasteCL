package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const FPASTE_URL = "http://fpaste.org"

type config struct {
	help   *bool
	priv   *bool
	user   *string
	pass   *string
	lang   *string
	expire *string
}

type fpaste struct {
}

var flags config

var filepaths []string

func init() {
	flag.BoolVar(flags.help, "help", "Display this help")
	flag.BoolVar(flags.priv, "private", "Private paste flag")
	flag.StringVar(flags.user, "user", "", "An alphanumeric username of the paste author")
	flag.StringVar(flags.pass, "pass", "", "Add a password")
	flag.StringVar(flags.lang, "lang", "Text", "The development language used")
	flag.StringVar(flags.expire, "expire", "0", "Seconds after which paste will be deleted from server")
	filepaths = flag.Args()
	flag.Parse()
}

func (f *fpaste) Get() (string, error) {

}

/*
Handling API errors so we can know why the request did not went as expetected.
The resquest will not be resend to prevent the user from being banned.
*/

func handleAPIError(error string) error {
	errorStr := make(map[string]string)
	errorStr["err_nothing_to_do"] = "No POST request was received by the create API"
	errorStr["err_author_numeric"] = "The paste author's alias should be alphanumeric"
	errorStr["err_save_error"] = "An error occurred while saving the paste"
	errorStr["err_spamguard_ipban"] = "Poster's IP address is banned"
	errorStr["err_spamguard_stealth"] = "The paste triggered the spam filter"
	errorStr["err_spamguard_noflood"] = "Poster is trying the flood"
	errorStr["err_spamguard_php"] = "Poster's IP address is listed as malicious"
	if err, ok := errorStr[error]; ok {
		return fmt.Errorf("API error: %s", err)
	}
	return fmt.Errorf("API error: Unknown [%s]", error)
}

func (f *fpaste) Put(s string) error {
	values := url.Values{
		"paste_data":     {string(src)},
		"paste_lang":     {*flags.lang},
		"api_submit":     {"true"},
		"mode":           {"json"},
		"paste_user":     {*flags.user},
		"paste_password": {*flags.pass},
	}
	if duration, err := time.ParseDuration(*flags.expire); err != nil {
		return err
	} else if secs := duration.Seconds(); secs >= 1 {
		values.Add("paste_expire", strconv.FormatFloat(secs, 'f', -1, 64))
	}
	if *flags.priv {
		values.Add("paste_private", "yes")
	}
	resp, erreq := http.PostForm(FPASTE_URL, values)
	if erreq != nil {
		return erreq
	}
	defer resp.Body.Close()
	type res struct {
		Id    string `json:"id"`
		Hash  string `json:"hash"`
		Error string `json:"error"`
	}
	type pasteUrls struct {
		Result res `json:"result"`
	}
	var m pasteUrls
	slice, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(slice, &m)
	if err != nil {
		return err
	}
	if m.Result.Error != "" {
		return handleAPIError(m.Result.Error)
	}
	fmt.Fprintf(os.Stdout, "%s/%s/%s\n", FPASTE_URL, m.Result.Id, m.Result.Hash)
	return nil
}
