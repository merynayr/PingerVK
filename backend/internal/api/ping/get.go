package ping

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/merynayr/PingerVK/pkg/sys/codes"
)

// Get - отправляет запрос в сервисный слой на получение данных о пингах
func (api *API) Get(ctx *gin.Context) {
	pingObj, err := api.pingService.Get(ctx)
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(int(codes.InternalServerError), gin.H{
			"error": fmt.Sprintf("Failed to fetch ping data: %v", err),
		})
		return
	}

	ctx.JSON(int(codes.OK), gin.H{
		"data": pingObj,
	})
}
