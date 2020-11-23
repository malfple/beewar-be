package message

import (
	"encoding/json"
	"errors"
)

const (
	// CmdShutdown is a cmd for shutting down game hub's listener
	CmdShutdown = "SHUTDOWN"
)

var (
	// ErrCmdNotAllowed is returned when client sends a restricted cmd
	ErrCmdNotAllowed = errors.New("that cmd is not allowed")
)

// GameMessageTemporary is the container struct when unmarshalling from json
type GameMessageTemporary struct {
	Cmd  string          `json:"cmd"`
	Data json.RawMessage `json:"data"`
}

// GameMessage is the main message struct for websocket message exchange
type GameMessage struct {
	Cmd    string      `json:"cmd"`
	Sender int64       `json:"sender,omitempty"`
	Data   interface{} `json:"data"`
}

// UnmarshalAndValidateGameMessage unmarshals the raw byte data into message struct
// also validates the cmd. if it is not allowed, it will return error
func UnmarshalAndValidateGameMessage(rawPayload []byte, senderID int64) (*GameMessage, error) {
	temp := &GameMessageTemporary{}
	err := json.Unmarshal(rawPayload, temp)
	if err != nil {
		return nil, err
	}

	message := &GameMessage{
		Cmd:    temp.Cmd,
		Sender: senderID,
	}

	switch temp.Cmd {
	case CmdShutdown:
		return nil, ErrCmdNotAllowed
	default:
		var data string
		err := json.Unmarshal(temp.Data, &data)
		if err != nil {
			return nil, err
		}
		message.Data = data
	}

	return message, nil
}

// MarshalGameMessage marshals game message into byte array
func MarshalGameMessage(message *GameMessage) ([]byte, error) {
	rawData, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}
