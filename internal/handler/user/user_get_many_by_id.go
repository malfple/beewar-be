package user

import (
	"encoding/json"
	"gitlab.com/beewar/beewar-be/internal/access"
	"gitlab.com/beewar/beewar-be/internal/access/model"
	"gitlab.com/beewar/beewar-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

// HandleUserGetManyByID handles getting many users by given id
func HandleUserGetManyByID(w http.ResponseWriter, r *http.Request) {
	ids := strings.Split(r.URL.Query().Get("ids"), ",")
	userIDs := make([]uint64, len(ids))
	for i, s := range ids {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		userIDs[i] = uint64(id)
	}

	users, err := access.QueryUsersByID(userIDs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := &GetManyByIDResponse{Users: users}
	for i := range resp.Users {
		if resp.Users[i] != nil {
			resp.Users[i].Password = "nopeeee"
		}
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// GetManyByIDResponse is a response for user get many by id handler
type GetManyByIDResponse struct {
	Users []*model.User `json:"users"`
}
