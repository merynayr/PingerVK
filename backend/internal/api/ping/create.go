package ping

import (
	"github.com/gin-gonic/gin"
	"github.com/merynayr/PingerVK/backend/internal/model"
	"github.com/merynayr/PingerVK/pkg/sys/codes"
)

// Create - отправляет запрос в сервисный слой на создание данных о пингах
func (api *API) Create(ctx *gin.Context) {
	var ping model.Pings

	if err := ctx.ShouldBindJSON(&ping); err != nil {
		_ = ctx.Error(err)
		ctx.JSON(int(codes.BadRequest), gin.H{"error": "invalid request"})
		return
	}

	id, err := api.pingService.Create(ctx, &ping)
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(int(codes.InternalServerError), gin.H{"error": "failed to create ping"})
		return
	}

	ctx.JSON(int(codes.OK), gin.H{"id": id})
}
