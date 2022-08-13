package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j))
	return []byte(stmp), nil
}

type UserResponse struct {
	ID       int32    `json:"id"`
	NikeName string   `json:"nike_name"`
	Mobile   string   `json:"mobile"`
	Birthday JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
}
