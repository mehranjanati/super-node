package domain

// User represents a user in the system.

type User struct {
	ID      int64  `json:"id"`
	Address string `json:"address"`
}

// AppData represents generic application data.

type AppData struct {
	ID   string `json:"id"`
	Data []byte `json:"data"`
}
