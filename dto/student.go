package dto

type StudentData struct {
	Name     string `json:"name" binding:"required,alphaspace,min=3,max=150"`
	Email    string `json:"email" binding:"required,email,max=254"`
	Class    string `json:"class" binding:"required,printascii"`
	Password string `json:"password" binding:"required,max=128"`
}

type StudentLogin struct {
	Email    string `json:"email" binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

type StudentResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Class string `json:"class"`
}
