package ping

import (
	"net/http"

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
	router.GET("/health", api.Health)
}

// Health проверка состояния
func (api *API) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
