package dto

type UserListResponse struct {
  ID       int     `json:"idUser"`
  Email    string  `json:"email"`
  FullName string  `json:"fullName"`
  Role     string  `json:"role"`
  Phone    *string `json:"phoneNumber,omitempty"`
  Picture  *string `json:"profilePicture,omitempty"`
}

type UpdateProfileRequest struct {
  FullName       *string `json:"fullName,omitempty"`
  PhoneNumber    *string `json:"phoneNumber,omitempty"`
  ProfilePicture *string `json:"profilePicture,omitempty"`
  OldPassword    *string `json:"oldPassword,omitempty"`
  NewPassword    *string `json:"newPassword,omitempty"`
}
