package models

type User struct {
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	Password    string       `json:"-"`
	Credentials []Credential `json:"credentials"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Credential struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
