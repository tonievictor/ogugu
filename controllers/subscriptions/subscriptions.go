package subscriptions

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"ogugu/services/subscriptions"
)

type SubscriptionController struct {
	log        *zap.Logger
	cache      *redis.Client
	subservice *subscriptions.SubscriptionService
}
