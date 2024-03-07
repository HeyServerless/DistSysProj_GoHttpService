package model

type ComputeRequest struct {
	Expressions []Expression `json:"expressions"`
	CallbackURL string       `json:"callback_url"`
}

type ExpressionRequest struct {
	Expression string `json:"expression" binding:"required"`
}

type ExpressionResultRequest struct {
	RequestId string `json:"request_id" binding:"required"`
}

type UpdateResponseRequest struct {
	RequestId string `json:"request_id" binding:"required"`
	Response  string `json:"response" binding:"required"`
}
