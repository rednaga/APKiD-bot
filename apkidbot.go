package main

import (
	"bytes"
	"encoding/hex"
	"crypto/sha256"
	"net/http"
	"io/ioutil"
	"fmt"
	"os"
	"strings"
	"os/exec"
	
	"github.com/nlopes/slack"
)

func main() {

	token := os.Getenv("SLACK_TOKEN")
	api := slack.New(token)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if(ev.SubType == "file_share" && (ev.File.Mimetype == "application/zip" || ev.File.Mimetype == "application/octet-stream" || ev.File.Mimetype == "application/vnd.android.package-archive")) {
					fmt.Printf("Found a file share! Attempting download...\n")
					respond(rtm, ev, prefix)
				}
				
				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					fmt.Printf("Found my prefix!\n")
					respond(rtm, ev, prefix)
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				//Take no action
			}
		}
	}
}

func respond(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	text := msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	if(msg.SubType == "file_share") {
		fmt.Printf("Detected fileshare, attempted to grab file : %s\n", msg.File.URLPrivate)
		text = msg.File.URLPrivate
	}

	
	hash, err := downloadFile(text)
	if err != nil {
		fmt.Printf("Error : ", err.Error())
		rtm.SendMessage(rtm.NewOutgoingMessage("Something bad happened with your input...", msg.Channel))
	} else {
		out := apkid("./data/" + hash)
		rtm.SendMessage(rtm.NewOutgoingMessage(out, msg.Channel))
	}
}

func downloadFile(url string) (filename string, err error) {
	// Get the data
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	// If url is a slack url - add the authrozation header
	// TODO: Extract this to a config / environment variable ?
	if(strings.HasPrefix(url, "https://rednaga.slack.com") || strings.HasPrefix(url, "https://files.slack.com")) {// || strings.HasPrefix(url, "https:\/\/files.slack.com")) {
		req.Header.Set("Authorization", "Bearer " + os.Getenv("SLACK_TOKEN"))
	}
	resp, err := client.Do(req)
	
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	
	// Hash data
	hasher := sha256.New()
	hasher.Write(data)
	hash := hex.EncodeToString(hasher.Sum(nil))

	fmt.Printf("Downloaded file from : %s\n", resp.Request.URL.String())

	// Check if file already exists
	if(exists(hash)) {
		return hash, nil
	}
	
	// Create the file
	out, err := os.Create("./data/" + hash)
	if err != nil  {
		return "", err
	}
	defer out.Close()

	
	_, err = out.Write(data)
	if err != nil  {
		return "", err
	}

	return hash, nil
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func apkid(name string) string {
	cmd := exec.Command("./APKiD/docker/apkid.sh", name)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return ""
	}
	fmt.Println("Result: " + out.String())

	return out.String()
}
