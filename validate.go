package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-product/common"
	"go-product/datamodels"
	"go-product/encrypt"
	"go-product/rabbitmq"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

// Set Cluster ip
var hostArray = []string{"127.0.0.1", "127.0.0.1"}

var localHost = ""

// Product number control inner server ip, Or getone SLB intranet server ip
var GetOneIp = "127.0.0.1"

var GetOnePort = "8084"

var port = "8083"

var hashConsistent *common.Consistent

// rabbitmq
var rabbitMqValidate *rabbitmq.RabbitMQ

// Save control information
type AccessControl struct {
	// save user want save info
	sourcesArray map[int]interface{}
	sync.RWMutex
}

// Create global variable
var accessControl = &AccessControl{sourcesArray: make(map[int]interface{})}

// Get customized data
func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.sourcesArray[uid]
	return data
}

// Set record
func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	m.sourcesArray[uid] = "Hello golang"
	m.RWMutex.Unlock()
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	// Get user UID
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}

	// Use hash algorithm judge server by uid
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	// check is local server
	if hostRequest == localHost {
		// execute data read
		return m.GetDataFromMap(uid.Value)
	} else {
		// proxy other server result
		return GetDataFromOtherMap(hostRequest, req)
	}
}

// Get localhost map result is boolean
func (m *AccessControl) GetDataFromMap(uid string) (isOk bool) {
	// @TODO Enable blow sections
	//uidInt, err := strconv.Atoi(uid)
	//if err != nil {
	//	return false
	//}
	//data := m.GetNewRecord(uidInt)
	//
	//if data != nil {
	//	return true
	//}
	//return false
	return true
}

// Get from other node
func GetDataFromOtherMap(host string, request *http.Request) bool {
	// proxy server mock visit
	hostUrl := "http://" + host + ":" + port + "/checkRight"
	response, body, err := GetCurl(hostUrl, request)
	if err != nil {
		return false
	}

	// Judge state
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}

	return false
}

// Simulate http request
func GetCurl(hostUrl string, request *http.Request) (response *http.Response, body []byte, err error) {
	// Get uid
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}

	// Get sign
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}

	// proxy server mock visit
	client := &http.Client{}
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		return
	}

	// Prepare cookie arguments
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}

	// Add cookie mock request
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	// Get response result
	response, err = client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(response.Body)
	return
}

func CheckRight(w http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
	}
	w.Write([]byte("true"))
	return
}

// Execute normal logic
func Check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Execute check!")
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"]) <= 0 {
		w.Write([]byte("false"))
		return
	}
	productString := queryForm["productID"][0]
	fmt.Println(productString)
	// Get user cookie
	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false"))
	}
	fmt.Println("User Cookie:" + userCookie.Value)

	// 1. Check distribute authorization
	right := accessControl.GetDistributedRight(r)
	if right == false {
		w.Write([]byte("false"))
	}

	// 2. Get number control prevent oversold in spikes
	hostUrl := "http://" + GetOneIp + ":" + GetOnePort + "/getOne"
	responseValidate, validateBody, err := GetCurl(hostUrl, r)
	if err != nil {
		w.Write([]byte("false"))
	}

	if responseValidate.StatusCode == 200 {
		if string(validateBody) == "true" {
			// arrangement order
			// 1. Get product id
			productID, err := strconv.ParseInt(productString, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}

			// 2. Get user ID
			userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}

			// 3. Create message structure
			message := datamodels.NewMessage(userID, productID)
			byteMessage, err := json.Marshal(message)

			if err != nil {
				w.Write([]byte("false"))
				return
			}

			// 4. Produce message
			err = rabbitMqValidate.PublishSimple(string(byteMessage))
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			w.Write([]byte("true"))
			return
		}
	}

	w.Write([]byte("false"))
	return
}

// User identify interceptor each every api
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Execute identify check!")
	// cookie authorization check
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserInfo(r *http.Request) error {
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("UID cookie get failed!")
	}

	// Get user sign
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("User sign Cookie get failed!")
	}

	// Compare information and decrypt
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("Sign has been changed!")
	}

	fmt.Println("------Compare result-------")
	fmt.Println("UID:" + uidCookie.Value)
	fmt.Println("Decrypt UID:" + string(signByte))

	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}

	return errors.New("validate user identify failed")
}

func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}

// Main entrance
func main() {
	hashConsistent = common.NewConsistent()
	// Use has consistent add node
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	localIp, err := common.GetInternalIp()
	if err != nil {
		fmt.Println(err)
	}

	localHost = localIp
	fmt.Println("Local IP:" + localHost)

	// Create rabbitmq instance
	rabbitMqValidate = rabbitmq.NewRabbitMQSimple("golang")
	defer rabbitMqValidate.Destory()

	// 1. filter
	filter := common.NewFilter()

	// Register interceptor
	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/checkRight", Auth)

	// 2. Start server
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight", filter.Handle(CheckRight))
	http.ListenAndServe(":8083", nil)

}
