package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
	Photo []Photo `json:"photo"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type Photo struct {
	FileId string `json:"file_id"`
	Width int `json:"width"`
	Height int `json:"height"`
}

type File struct {
	FileId string `json:"file_id"`
	FilePath string `json:"file_path"`
}

type FileResponse struct {
	Ok     bool     `json:"ok"`
	Result File `json:"result"`}
