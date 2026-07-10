package console

type Result struct {
	OK        bool   `json:"ok"`
	OutDir    string `json:"out_dir,omitempty"`
	PageCount int    `json:"page_count,omitempty"`
	Error     string `json:"error,omitempty"`
}