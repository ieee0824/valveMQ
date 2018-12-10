package valve

import (
	"fmt"
	"time"
)

type limit uint

func (l limit) String() string {
	return fmt.Sprintf("%d dq per sec", l)
}

func (l limit) DqSpan() time.Duration {
	if l == 0 {
		return 0
	}
	return time.Duration(float64(time.Second) / float64(l))
}

type Setting struct {
	Limit     limit  `json:"limit" sql:"limit"`
	QueueName string `json:"queue_name" sql:"queue_name"`
}
