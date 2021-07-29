package campaign

import "github.com/gorilla/mux"

// RegisterCampaignRouter builds router for campaign
func RegisterCampaignRouter(router *mux.Router) {
	router.HandleFunc("/start", HandleStartCampaign).Methods("POST")
}
