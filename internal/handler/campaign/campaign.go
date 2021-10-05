package campaign

import "github.com/gorilla/mux"

// RegisterCampaignRouter builds router for campaign
func RegisterCampaignRouter(router *mux.Router) {
	router.HandleFunc("/list", HandleCampaignList).Methods("GET")
	router.HandleFunc("/current", HandleCurrentCampaign).Methods("GET")
	router.HandleFunc("/start", HandleStartCampaign).Methods("POST")
}
