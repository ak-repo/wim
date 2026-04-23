package model

type ProductCategoryDTO struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	RefCode  string `db:"ref_code"`
	IsActive bool   `db:"is_active"`
}

type ProductCategoryResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	RefCode  string `json:"refCode"`
	IsActive bool   `json:"isActive"`
}

type ProductCategoryRequest struct {
	ID       int    `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	RefCode  string `json:"refCode,omitempty"`
	IsActive *bool  `json:"isActive,omitempty"`
}

type ProductCategoryParams struct {
	Active *bool `json:"active"`
	Page   int   `json:"page"`
	Limit  int   `json:"limit"`
}

// API conversion
type ProductCategoryDTOs []*ProductCategoryDTO

func (m *ProductCategoryDTOs) ToAPIResponse() []*ProductCategoryResponse {
	var responses []*ProductCategoryResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPIResponse())
	}
	return responses
}

func (m *ProductCategoryDTO) ToAPIResponse() *ProductCategoryResponse {
	return &ProductCategoryResponse{
		ID:       m.ID,
		Name:     m.Name,
		RefCode:  m.RefCode,
		IsActive: m.IsActive,
	}
}
