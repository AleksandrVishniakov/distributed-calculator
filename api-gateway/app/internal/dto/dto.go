package dto

type CalculationRequestDTO struct {
	Expression string `json:"expression"`
}

type CalculationResponseDTO struct {
	Id     int    `json:"id"`
	Result string `json:"result"`
}

type LimitOffsetRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type OperationsSetup struct {
	AddTime      int64 `json:"addTime"`
	SubtractTime int64 `json:"subtractTime"`
	MultiplyTime int64 `json:"multiplyTime"`
	DivideTime   int64 `json:"divideTime"`
}
