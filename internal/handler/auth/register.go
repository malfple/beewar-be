package auth

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/controller/auth"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

// HandleRegister handles new user registration
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req := &RegisterRequest{}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errMsg := ""
	if err := auth.Register(req.Email, req.Username, req.Password); err != nil {
		errMsg = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := &RegisterResponse{
		ErrMsg: errMsg,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// RegisterRequest is a request struct
type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterResponse is a response for register handler
// if err_msg is empty string, success
type RegisterResponse struct {
	ErrMsg string `json:"err_msg"`
}
