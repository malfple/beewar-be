package message

import (
	"encoding/json"
	"errors"
)

// some Cmd is restricted to client-sent only or server-sent only
// let's just simplify these as client-only and server-only
const (
	// CmdShutdown is a cmd for shutting down game hub's listener. server-only
	CmdShutdown = "SHUTDOWN"
	// CmdChat is a cmd for text chat. no restriction
	CmdChat = "CHAT"
	// CmdGameInfo is a cmd for sending game data. server-only
	CmdGameData = "GAME_DATA"
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
	Sender uint64      `json:"sender,omitempty"`
	Data   interface{} `json:"data"`
}

// UnmarshalAndValidateGameMessage unmarshals the raw byte data into message struct
// also validates the cmd. if it is not allowed for client, it will return error
func UnmarshalAndValidateGameMessage(rawPayload []byte, senderID uint64) (*GameMessage, error) {
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
	case CmdChat:
		var data string
		err := json.Unmarshal(temp.Data, &data)
		if err != nil {
			return nil, err
		}
		message.Data = data
	default:
		return nil, ErrCmdNotAllowed
	}

	return message, nil
}

// MarshalGameMessage marshals game message into byte array
// this does not verify whether a cmd is restricted for server
// we thrust whoever sends the message to client is responsible
func MarshalGameMessage(message *GameMessage) ([]byte, error) {
	rawData, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}
