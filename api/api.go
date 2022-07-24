package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"net/smtp"
)

var emails []string

func getRate() (string, int) {
	var resp string
	var statusCode int
	respone, err := http.Get("https://api.coinbase.com/v2/prices/spot?currency=UAH")

	if err != nil {
		statusCode = http.StatusBadRequest
		resp = "Invalid status value"
	} else {
		statusCode = http.StatusOK
		data, _ := ioutil.ReadAll(respone.Body)

		var jsonMap map[string]map[string]string
		err := json.Unmarshal(data, &jsonMap)
		if err != nil {
			statusCode = http.StatusBadRequest
			resp = "Error happened in respond JSON marshal"
		} else {
			statusCode = http.StatusOK
			resp = jsonMap["data"]["amount"]
		}
	}

	return resp, statusCode
}

func readEmails(filename string) {

	file, err := os.Open("emails.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	emails = nil
	for scanner.Scan() {
		emails = append(emails, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func sendEmails() {
	from := "sender.email@example.com"
	user := "d6715ed2eb5bc7"
	password := "8490b88643af41"

	to := emails

	addr := "smtp.mailtrap.io:2525"
    host := "smtp.mailtrap.io"

	body, _ := getRate()

	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += "Subject: BTC to UAH rate\r\n"
	message += fmt.Sprintf("\r\n%s\r\n", body)

	auth := smtp.PlainAuth("", user, password, host)
	err := smtp.SendMail(addr, auth, from, to, []byte(message))
	if err != nil {
	    panic(err)
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func addEmail(filename string, email string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Could not open " + filename)
		return
	}
	defer file.Close()

	_, err2 := file.WriteString(email + "\n")
	if err2 != nil {
		fmt.Println("Could not write email")
		return
	}

	emails = append(emails, email)
}

func handleRequests(w http.ResponseWriter, r *http.Request) {

	resp := make(map[string]string)
	var statusCode int

	switch r.Method {
	case "GET":

		switch r.URL.Path {
		case "/rate":
			resp["message"], statusCode = getRate()

		default:
			return
		}

	case "POST":

		switch r.URL.Path {
		case "/subscribe":
			r.ParseForm()
			email := r.Form["email"][0]

			if contains(emails, email) {
				statusCode = http.StatusConflict
				resp["message"] = "E-mail вже існує"
			} else {
				addEmail("Emails.txt", email)

				statusCode = http.StatusOK
				resp["message"] = "E-mail додано"
			}

		case "/sendEmails":
			sendEmails()

			statusCode = http.StatusOK
			resp["message"] = "E-mailʼи відправлено"

		default:
			fmt.Println("Undefined path")
		}

	default:
		return
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in respond JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

}

func StartServer(port string) error {
	readEmails("Emails.txt")

	http.HandleFunc("/", handleRequests)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return fmt.Errorf("Could not start client API server on port %s: %w", port, err)
	}

	return nil
}