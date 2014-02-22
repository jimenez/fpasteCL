package main

import (
	getopt "code.google.com/p/getopt"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	http "net/http"
	url "net/url"
	"os"
	"strconv"
)

const FPASTE_URL = "http://fpaste.org"

func main() {
	opts := initConfig(os.Args)
	if *opts.help {
		getopt.PrintUsage(os.Stdout)
	} else {
		/*Passing stdin for test mocking*/
		files, errs := handleArgs(os.Stdin, getopt.CommandLine)
		for _, file := range files {
			if len(file) != 0 {
				if err := copyPaste(file, opts); err != nil {
					errs = append(errs, err)
				}
			}
		}
		for _, err := range errs {
			log.Print(err)
		}
		if len(errs) > 0 {
			os.Exit(-1)
		}
	}
}

type config struct {
	help   *bool
	priv   *bool
	user   *string
	pass   *string
	lang   *string
	expire *int
}

func initConfig(args []string) *config {
	getopt.CommandLine = getopt.New()
	var flags config
	flags.help = getopt.BoolLong("help", 'h', "Display this help")
	flags.priv = getopt.BoolLong("private", 'P', "Private paste flag")
	flags.user = getopt.StringLong("user", 'u', "", "An alphanumeric username of the paste author")
	flags.pass = getopt.StringLong("pass", 'p', "", "Add a password")
	flags.lang = getopt.StringLong("lang", 'l', "Text", "The development language used")
	flags.expire = getopt.IntLong("expire", 'e', 0, "Seconds after which paste will be deleted from server")
	getopt.SetParameters("[FILE...]")
	getopt.CommandLine.Parse(args)
	return &flags
}

func handleArgs(stdin io.Reader, commandLine *getopt.Set) (files [][]byte, errs []error) {
	if commandLine.NArgs() > 0 {
		for _, x := range commandLine.Args() {
			file, err := os.Open(x)
			if err != nil {
				errs = append(errs, fmt.Errorf("Skipping [FILE: %s] since it cannot be opened (%s)", x, err))
				continue
			}
			defer file.Close()
			data, erread := ioutil.ReadAll(file)
			if erread != nil {
				errs = append(errs, fmt.Errorf("Skipping [FILE: %s] since it cannot be read (%s)", x, erread))
			} else {
				files = append(files, data)
			}
		}
	} else {
		data, erread := ioutil.ReadAll(stdin)
		if erread != nil {
			errs = append(errs, erread)
		} else {
			files = append(files, data)
		}
	}
	return files, errs
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

func copyPaste(src []byte, opts *config) error {
	values := url.Values{
		"paste_data":     {string(src)},
		"paste_lang":     {*opts.lang},
		"api_submit":     {"true"},
		"mode":           {"json"},
		"paste_user":     {*opts.user},
		"paste_password": {*opts.pass},
		"paste_expire":   {strconv.Itoa(*opts.expire)},
	}
	if *opts.priv {
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
