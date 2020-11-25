package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalGameMessage(t *testing.T) {
	rawMessage := []byte(`{"cmd":"CHAT","data":"hellow"}`)
	msg, err := UnmarshalAndValidateGameMessage(rawMessage, 123)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(123), msg.Sender)
	assert.Equal(t, "CHAT", msg.Cmd)
	assert.Equal(t, "hellow", msg.Data)
}

func TestMarshalGameMessage(t *testing.T) {
	rawMessage := []byte(`{"cmd":"RANDOM_CMD","data":"hellow"}`)
	msg := &GameMessage{
		Cmd:  "RANDOM_CMD",
		Data: "hellow",
	}
	rawMessageExpected, err := MarshalGameMessage(msg)
	assert.Equal(t, nil, err)
	assert.Equal(t, rawMessage, rawMessageExpected)
}

func TestMarshalGameMessage2(t *testing.T) {
	rawMessage := []byte(`{"cmd":"SHUTDOWN","data":null}`)
	msg := &GameMessage{
		Cmd: "SHUTDOWN",
	}
	rawMessageExpected, err := MarshalGameMessage(msg)
	assert.Equal(t, nil, err)
	assert.Equal(t, rawMessage, rawMessageExpected)
}
