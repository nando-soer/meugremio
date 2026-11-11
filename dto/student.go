package dto

type AvatarData struct {
	Body      int
	Shirt     int
	Pants     int
	Accessory int
}

type ClassData struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type StudentData struct {
	Name        string     `json:"name" binding:"required,min=3,max=150"`
	Email       string     `json:"email" binding:"required,email,max=254"`
	Class       string     `json:"class"`
	Description string     `json:"description"`
	Banner      int        `json:"banner"`
	Avatar      AvatarData `json:"avatar"`
	Password    string     `json:"password" binding:"required,max=128"`
}

type StudentEditData struct {
	Name        string     `json:"name"`
	Class       string     `json:"class"`
	Description string     `json:"description"`
	Banner      int        `json:"banner"`
	Avatar      AvatarData `json:"avatar"`
}

type StudentLogin struct {
	Email    string `json:"email" binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

type StudentResponse struct {
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Class       string     `json:"class"`
	Description string     `json:"description"`
	Banner      int        `json:"banner"`
	Avatar      AvatarData `json:"avatar"`
}

type StudentPasswordUpdate struct {
	Token    string `json:"token"`
	Password string `json:"password,omitempty"`
}
