package pandora_go

import "testing"

func TestErrorResponseUnMarshal(t *testing.T) {
	datas := [][]byte{
		[]byte(`{"detail":{"message":"You have sent too many messages to the model. Please try again later.","code":"model_cap_exceeded","clears_in":9665}}`),
		[]byte(`{"detail":"Only one message at a time. Please allow any other responses to complete before sending another message, or wait one minute."}`),
	}
	for _, v := range datas {
		result, err := NewErrorResponse(v)
		if err != nil {
			t.Error(err)
			continue
		}
		t.Logf("\nErrorResponse: %#v\nErrorResponse.Detail.GetType(): %#v\nErrorResponse.Detail: %#v\n", result, result.Detail.GetType(), result.Detail)
	}
}
