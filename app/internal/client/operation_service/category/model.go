package category

type Category struct {
	UUID     string `json:"uuid"`
	UserUUID string `json:"user_uuid"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

type CreateCategoryDTO struct {
	UserUUID string `json:"user_uuid"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

type UpdateCategoryDTO struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
