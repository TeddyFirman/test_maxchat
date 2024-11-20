package models

type Data struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Model       string `json:"model"`
	Tech        string `json:"tech"`
	Status      string `json:"status"`
	Description string `json:"description,omitempty"` // Optional
}

type ModelReference struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TechReference struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
