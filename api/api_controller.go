package api

import (
	"belli/onki-game-ideas-mongo-backend/datastore"
	"belli/onki-game-ideas-mongo-backend/service/cache"
)

type APIController struct {
	access datastore.Access
	cache  cache.Service
}

func NewAPIController(access datastore.Access, RedisService *cache.RedisService) *APIController {
	return &APIController{
		access: access,
		cache:  RedisService,
	}
}
