package users

// NewUser indicates that the given user has connected.
type NewUser struct {
	// Cid is the user's connection ID
	Cid int

	// SendMessage is a callable to send a message to this user (with blocking
	// semantics similar to RequestAsync)
	SendMessage func(string)
}

// UserGone indicates that the given user has disconnected.
type UserGone struct {
	// Cid is the user's connection ID
	Cid int
}

// UserMessage carries a message from a user.
type UserMessage struct {
	// Cid is the user's connection ID
	Cid int
	// Message is the message (without trailing newline)
	Message string
}

type user struct {
	cid         int
	room        string
	sendMessage func(string)
}
