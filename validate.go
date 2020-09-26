package main

import (
	"errors"
	"fmt"
	"go-product/common"
	"go-product/encrypt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

// Set Cluster ip
var hostArray = []string{"127.0.0.1", "192.168.56.38"}

var localHost = "192.168.56.38"

var port = "8081"

var hashConsistent *common.Consistent

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
	m.sourcesArray[uid] = "Hello goper"
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
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	data := m.GetNewRecord(uidInt)

	if data != nil {
		return true
	}
	return false
}

// Get from other node
func GetDataFromOtherMap(host string, request *http.Request) bool {
	// Get uid
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return false
	}

	// Get sign
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return false
	}

	// proxy server mock visite
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+host+":"+port+"/check", nil)
	if err != nil {
		return false
	}

	// Prepare cookie arguments
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}

	// Add cookie mock request
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	// Get response result
	response, err := client.Do(req)
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(response.Body)

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

// Execute normal logic
func Check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Execute check!")
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

	// 1. filter
	filter := common.NewFilter()

	// Register interceptor
	filter.RegisterFilterUri("/check", Auth)

	// 2. Start server
	http.HandleFunc("/check", filter.Handle(Check))
	http.ListenAndServe(":8083", nil)

}
