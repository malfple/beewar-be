package message

import "encoding/json"

// GameMessageTemporary is the container struct when unmarshalling from json
type GameMessageTemporary struct {
	Cmd  string          `json:"cmd"`
	Data json.RawMessage `json:"data"`
}

// GameMessage is the main message struct for websocket message exchange
type GameMessage struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

// UnmarshalGameMessage unmarshals the raw byte data into message struct
func UnmarshalGameMessage(rawPayload []byte) (*GameMessage, error) {
	temp := &GameMessageTemporary{}
	err := json.Unmarshal(rawPayload, temp)
	if err != nil {
		return nil, err
	}

	message := &GameMessage{
		Cmd: temp.Cmd,
	}

	switch temp.Cmd {
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

// MarshalGameMessage
func MarshalGameMessage(message *GameMessage) ([]byte, error) {
	rawData, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}
