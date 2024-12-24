package api

type Message struct {
	Room    string
	Content []byte
}

type BroadcastRequest struct {
	Message Message
	Sender  *Client
}
