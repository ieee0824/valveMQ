package valve

type Setting struct {
	LimitMPS  float64 `json:"limit_mps" sql:"limit_mps"`
	QueueName string  `json:"queue_name" sql:"queue_name"`
}
