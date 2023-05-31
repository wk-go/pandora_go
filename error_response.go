package pandora_go

import (
	"encoding/json"
)

type ErrorInterface interface {
	GetType() string
	GetCode() string
	GetMessage() string
	GetClearsIn() int
}

type ErrorJson struct {
	Message  string `json:"message"`
	Code     string `json:"code"`
	ClearsIn int    `json:"clears_in"`
}

func (e *ErrorJson) GetType() string {
	return "json"
}

func (e *ErrorJson) GetCode() string {
	return e.Code
}
func (e *ErrorJson) GetMessage() string {
	return e.Message
}

func (e *ErrorJson) GetClearsIn() int {
	return e.ClearsIn
}

type ErrorString string

func (e *ErrorString) GetType() string {
	return "string"
}

func (e *ErrorString) GetCode() string {
	return ""
}
func (e *ErrorString) GetMessage() string {
	return string(*e)
}

func (e *ErrorString) GetClearsIn() int {
	return 0
}

// ErrorResponse 错误响应
// 响应示例1: {"detail":{"message":"You have sent too many messages to the model. Please try again later.","code":"model_cap_exceeded","clears_in":9665}}
// 响应示例2: {"detail":"Only one message at a time. Please allow any other responses to complete before sending another message, or wait one minute."}
type ErrorResponse struct {
	Detail ErrorInterface `json:"detail"`
}

func NewErrorResponse(data ...[]byte) (result *ErrorResponse, err error) {
	result = new(ErrorResponse)
	if len(data) == 0 {
		return result, nil
	}
	err = json.Unmarshal(data[0], result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (er *ErrorResponse) UnmarshalJSON(data []byte) (err error) {
	tmp := map[string]any{}
	//result.Detail = ErrorJson{}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	var detail ErrorInterface
	detailBytes, _ := json.Marshal(tmp["detail"])

	switch tmp["detail"].(type) {
	case string:
		detail = new(ErrorString)
	case map[string]any:
		detail = new(ErrorJson)
	}

	err = json.Unmarshal(detailBytes, &detail)
	if err != nil {
		return err
	}
	er.Detail = detail
	return
}
