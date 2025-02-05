package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/merynayr/PingerVK/backend/internal/service"
)

// API ping структура
// объект сервисного слоя (его интерфейса)
type API struct {
	pingService service.PingService
}

// NewAPI возвращает новый объект имплементации API-слоя
func NewAPI(pingService service.PingService) *API {
	return &API{
		pingService: pingService,
	}
}

// RegisterRoutes регистрирует маршруты
func (api *API) RegisterRoutes(router *gin.Engine) {
	router.GET("/ping", api.Get)
	router.POST("/ping", api.Create)
}
