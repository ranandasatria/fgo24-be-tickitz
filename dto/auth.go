package dto


type AuthRegisterLogin struct {
	Email    string `form:"email" json:"email" db:"email" binding:"required,email"`
	Password string `form:"password" json:"password" db:"password" binding:"required"`
}

type UserResponse struct {
  ID       int    `json:"idUser"`
  Email    string `json:"email"`
  FullName string `json:"fullName"`
  Role     string `json:"role"`
}
