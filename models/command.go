package models

const (
	Subscribe      = "SUBSCRIBE"
	Unsubscribe    = "UNSUBSCRIBE"
	NumConnections = "NUM_CONNECTIONS"
)

type CommandBody struct {
	Command string `json:"command"`
}
