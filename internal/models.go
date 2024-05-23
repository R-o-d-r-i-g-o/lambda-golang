package internal

import "encoding/json"

const (
	MAX_NUM_OF_MESSAGES = 10
	WAIT_TIME_SECONDS   = 1
)

// Message holds the base structure of message.
type Message struct {
	Sender Sender     `json:"sender"`
	Order  Order[any] `json:"order"`
}

// Sender login witch notifies order.
type Sender struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Order holds request data.
type Order[T any] struct {
	Detail   OrderDetail `json:"detail"`
	Metadata T           `json:"metadata"`
}

// OrderDetail holds customer data.
type OrderDetail struct {
	Phone    string `json:"phone"`
	Customer string `json:"customer"`
}

// unmarshalResponse parses json response into a formated struct
func unmarshalResponse[T any](jsonMsg string) (res *T, err error) {
	return res,
		json.Unmarshal([]byte(jsonMsg), res)
}

// marshalElement encrypts response into josn syntax.
func marshalElement(el any) ([]byte, error) {
	return json.Marshal(el)
}
