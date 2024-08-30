package dto

type UserProfileRequest struct {
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
}

type UserRequest struct {
	Username string  `json:"username" validate:"required,min=3,max=32"`
	Email    string  `json:"email" validate:"required,email"`
	Name     string  `json:"name" validate:"min=3,max=32" validate:"required"`
	Phone    *string `json:"phone"`
	Photo    *string `json:"photo"`
	Status   string  `json:"status"`
	Role     string  `json:"role"`
	Address  *string `json:"address,omitempty"`
}
