package models

type InitPasswordResetReqBody struct {
	Email string `json:"email" binding:"required" example:"john.doe@example.com"`
}

func NewInitPasswordResetReqBody(email string) InitPasswordResetReqBody {
	return InitPasswordResetReqBody{
		Email: email,
	}
}

type ResetPasswordReqBody struct {
	Password           string `json:"password" binding:"required,min=16" example:"qwerty!123456789"` // TODO: Improve password strength
	ConfirmPassword    string `json:"confirmPassword" binding:"required,eqfield=Password" example:"qwerty!123456789"`
	PasswordResetToken string `json:"passwordResetToken" binding:"required" example:"ec6c66fbd3d92b1ad44f21613c5ee2e82c3dd65e8c918945308087ce77b5fe47"`
}

func NewResetPasswordReqBody(password, confirmPassword, passwordResetToken string) ResetPasswordReqBody {
	return ResetPasswordReqBody{
		Password:           password,
		ConfirmPassword:    confirmPassword,
		PasswordResetToken: passwordResetToken,
	}
}
