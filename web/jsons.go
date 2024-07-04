package web

type JsonID struct {
	Zone string `form:"zone"`
}

type JsonCommand struct {
	Zone    string `json:"zone" binding:"required"`
	Command string `json:"command" binding:"required"`
}
