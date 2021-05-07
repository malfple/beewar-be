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
	// CmdGameData is a cmd for sending game data. no restriction.
	// If sent by client, it means client is requesting game data.
	CmdGameData = "GAME_DATA"
	// CmdJoin is a cmd for joining game. no restriction
	CmdJoin = "JOIN"
	// CmdPing is a cmd for regular ping to avoid disconnection, server-only
	CmdPing = "PING"
	// CmdError is a cmd for sending server error. server-only
	CmdError = "ERROR"

	// CmdUnitMove is a unit cmd for general moving. no restriction.
	CmdUnitMove = "UNIT_MOVE"
	// CmdUnitAttack is a unit cmd for general attack. no restriction
	CmdUnitAttack = "UNIT_ATTACK"
	// CmdUnitMoveAndAttack is a unit cmd for general move and attack. no restriction.
	CmdUnitMoveAndAttack = "UNIT_MOVE_ATTACK"

	// CmdEndTurn is a cmd to end the turn
	CmdEndTurn = "END_TURN"
)

var (
	// ErrCmdNotAllowed is returned when client sends a restricted cmd
	ErrCmdNotAllowed = errors.New("message validation: that cmd is not allowed")
)

// GameMessageTemporary is the container struct when unmarshalling from json
type GameMessageTemporary struct {
	Cmd  string          `json:"cmd"`
	Data json.RawMessage `json:"data"`
}

// GameMessage is the main message struct for websocket message exchange
// WARNING: received messages should not be modified!
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
		if err := json.Unmarshal(temp.Data, &data); err != nil {
			return nil, err
		}
		message.Data = data
	case CmdGameData:
		// do nothing
	case CmdJoin:
		var data *JoinMessageData
		if err := json.Unmarshal(temp.Data, &data); err != nil {
			return nil, err
		}
		message.Data = data
	case CmdPing:
		return nil, ErrCmdNotAllowed
	case CmdUnitMove:
		var data *UnitMoveMessageData
		if err := json.Unmarshal(temp.Data, &data); err != nil {
			return nil, err
		}
		message.Data = data
	case CmdUnitAttack:
		var data *UnitAttackMessageData
		if err := json.Unmarshal(temp.Data, &data); err != nil {
			return nil, err
		}
		message.Data = data
	case CmdUnitMoveAndAttack:
		var data *UnitMoveAndAttackMessageData
		if err := json.Unmarshal(temp.Data, &data); err != nil {
			return nil, err
		}
		message.Data = data
	case CmdEndTurn:
		// do nothing
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

// GameErrorMessage creates an error game message
func GameErrorMessage(errMsg string) *GameMessage {
	return &GameMessage{
		Cmd:  CmdError,
		Data: errMsg,
	}
}
