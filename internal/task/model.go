package task

type Task struct {
	ID      string   `json:"id"`
	Status  string   `json:"status"`
	Files   []File   `json:"files"`
	Errors  []string `json:"errors,omitempty"`
	Archive string   `json:"archive,omitempty"`
}

type File struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}
