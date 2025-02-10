package main

import (
	"net/http"
	"fmt"
	"github.com/krishna-shrivastav232702/rssAggregator/internal/database"
	"github.com/krishna-shrivastav232702/rssAggregator/internal/auth"
)

type autheHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler autheHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldnt get user %v", err))
			return
		}
		handler(w,r,user)
	}
}
