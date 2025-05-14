package web

type JsonID struct {
	Zone string `form:"zone"`
	ID   string `json:"id" binding:"required"`
}

type JsonCommand struct {
	Zone    string `json:"zone" binding:"required"`
	Command string `json:"command" binding:"required"`
	ID      string `json:"id" binding:"required"`
}
