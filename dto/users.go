package dto

type UserListResponse struct {
  ID       int     `json:"idUser"`
  Email    string  `json:"email"`
  FullName string  `json:"fullName"`
  Role     string  `json:"role"`
  Phone    *string `json:"phoneNumber,omitempty"`
  Picture  *string `json:"profilePicture,omitempty"`
}
