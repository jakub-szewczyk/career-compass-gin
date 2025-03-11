package models

type InitPasswordResetReqBody struct {
	Email string `json:"email" binding:"required" example:"john.doe@example.com"`
}

func NewInitPasswordResetReqBody(email string) InitPasswordResetReqBody {
	return InitPasswordResetReqBody{
		Email: email,
	}
}
