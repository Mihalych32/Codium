package entity

type ExecuteRequest struct {
	Content  string `json:"content"`
	LangSlug string `json:"lang_slug"`
}
