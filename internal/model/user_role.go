package model

type UserRoleDTO struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	RefCode  string `db:"ref_code"`
	IsActive bool   `db:"is_active"`
}

type UserRoleResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	RefCode  string `json:"refCode"`
	IsActive bool   `json:"isActive"`
}

type UserRoleRequest struct {
	ID       int    `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	RefCode  string `json:"refCode,omitempty"`
	IsActive *bool  `json:"isActive,omitempty"`
}

type UserRoleParams struct {
	Active *bool `json:"active"`
	Page   int   `json:"page"`
	Limit  int   `json:"limit"`
}

// API conversion
type UserRoleDTOs []*UserRoleDTO

func (m *UserRoleDTOs) ToAPIResponse() []*UserRoleResponse {
	var responses []*UserRoleResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPIResponse())
	}
	return responses
}

func (m *UserRoleDTO) ToAPIResponse() *UserRoleResponse {
	return &UserRoleResponse{
		ID:       m.ID,
		Name:     m.Name,
		RefCode:  m.RefCode,
		IsActive: m.IsActive,
	}
}
