package oauth2

import (
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"gitlab.com/otqee/otqee-be/internal/logger"
	"net/http"
)

var manager *manage.Manager
var clientStore *store.ClientStore
var oauth2Server *server.Server

// InitOAuth2Server init oauth2Server
func InitOAuth2Server() {
	logger.Logger.Info("init oauth2Server")

	manager = manage.NewDefaultManager()
	// token memory storage
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	// client memory storage
	clientStore = store.NewClientStore()

	// TODO: remove
	clientStore.Set("111", &models.Client{
		ID:     "111",
		Secret: "222",
		Domain: "wtfdomain",
	})

	manager.MapClientStorage(clientStore)

	oauth2Server = server.NewDefaultServer(manager)
	oauth2Server.SetAllowGetAccessRequest(true)
	oauth2Server.SetClientInfoHandler(server.ClientFormHandler)

	//oauth2Server.SetInternalErrorHandler(func(err error) (re *errors.Response) {
	//	logger.Logger.Error("oauth2 internal error", zap.Error(err))
	//	return
	//})
	//
	//oauth2Server.SetResponseErrorHandler(func(re *errors.Response) {
	//	logger.Logger.Error("oauth2 response error", zap.Error(re.Error))
	//})
}

// HandleAuthorize is used to handle authorize request
func HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	err := oauth2Server.HandleAuthorizeRequest(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// HandleToken is used to handle token request
func HandleToken(w http.ResponseWriter, r *http.Request) {
	oauth2Server.HandleTokenRequest(w, r)
}
