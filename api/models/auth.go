package models

type SignUpReqBody struct {
	FirstName       string `json:"firstName" binding:"required" example:"John"`
	LastName        string `json:"lastName" binding:"required" example:"Doe"`
	Email           string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password        string `json:"password" binding:"required,min=16" example:"qwerty!123456789"` // TODO: Improve password strength
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password" example:"qwerty!123456789"`
}

func NewSignUpReqBody(firstName, lastName, email, password, confirmPassword string) SignUpReqBody {
	return SignUpReqBody{
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		Password:        password,
		ConfirmPassword: confirmPassword,
	}
}

type SignUpResBody struct {
	User  ProfileResBody `json:"user"`
	Token string         `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk4MDQ1NTEsInN1YiI6ImpvaG4uZG9lQGdtYWlsLmNvbSIsInVpZCI6IjZiZTA1YTcyLTc5OGQtNGI3Ny1iOGQzLTc3MjNhN2JmM2FkYSJ9.5sj2fHB3pky3N6-mDgaPQCQA0gkEz4oQsdtVEC9BLqE"`
}

func NewSignUpResBody(user AnyUser, token string) (*SignUpResBody, error) {
	profileResBody, err := NewProfileResBody(user)
	if err != nil {
		return nil, err
	}

	return &SignUpResBody{
		User:  *profileResBody,
		Token: token,
	}, nil
}

type SignInReqBody struct {
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" binding:"required,min=16" example:"qwerty!123456789"` // TODO: Improve password strength
}

func NewSignInReqBody(email, password string) SignInReqBody {
	return SignInReqBody{
		Email:    email,
		Password: password,
	}
}

type SignInResBody struct {
	User  ProfileResBody `json:"user"`
	Token string         `json:"token"`
}

func NewSignInResBody(user AnyUser, token string) (*SignInResBody, error) {
	profileResBody, err := NewProfileResBody(user)
	if err != nil {
		return nil, err
	}

	return &SignInResBody{
		User:  *profileResBody,
		Token: token,
	}, nil
}
