package Model

import (
	"github.com/go-xorm/xorm"
	json2 "encoding/json"
	"log"
)

var (
	Error_UserNotExsit = 1
	Error_UserAlreadyExsit = 2
	Error_ParamsNotCorrect = 3
	Error_PassowrdNotCorrect = 4
	Error_OperationFailed = 5
	Error_NotToken = 6
)

var x *xorm.Engine

func MakeJsonSuccess(success bool, message string) map[string]interface{}  {
	m := make(map[string]interface{})
	m["success"] = success
	m["message"] = message
	return m
}
func MakeJsonData(success bool, message string, data interface{}) map[string]interface{}  {
	m := make(map[string]interface{})
	m["success"] = success
	m["message"] = message
	m["data"]    = data
	return m
}
func MakeJsonError(errorCode int, message string) map[string]interface{} {
	m := make(map[string]interface{})
	m["errorCode"] = errorCode
	m["message"]   = message
	return m
}

func MakeJson(obj interface{}) interface{} {
	json, err := json2.Marshal(obj)
	if err != nil {
		log.Fatalf("MakeJson:%v\n", err)
		return nil
	}
	return json
}

func init() {

}

type LoginSuccessToken struct {
	Message string `json:"message"`
	Success bool `json:"success"`
	Data Token `json:"data"`
}


func MakeJsonForLoginSuccess(success bool, message string, data Token) LoginSuccessToken  {
	return LoginSuccessToken{Message:message, Success:success, Data:data}
}