package mymap

import (
	"encoding/json"
	"gitlab.com/otqee/otqee-be/internal/access"
	"gitlab.com/otqee/otqee-be/internal/access/model"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// HandleMapGet handles single map request
func HandleMapGet(w http.ResponseWriter, r *http.Request) {
	resp := &MapGetResponse{}

	mapID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err == nil { // id has to be integer
		mapp := access.QueryMapByID(uint64(mapID))
		resp.Map = mapp
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.GetLogger().Error("error encode", zap.Error(err))
	}
}

// MapGetResponse is a response for map get handler
type MapGetResponse struct {
	Map *model.Map `json:"map"`
}
