package web

type JsonNoID struct {
	Command string `json:"command" binding:"required"`
}

type JsonID struct {
	Zone string `json:"zone" binding:"required"`
	ID   string `json:"id" binding:"required"`
}

type JsonCommand struct {
	Zone    string `json:"zone" binding:"required"`
	Command string `json:"command" binding:"required"`
	ID      string `json:"id" binding:"required"`
}
