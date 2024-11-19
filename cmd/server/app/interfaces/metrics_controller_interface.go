package interfaces

import "github.com/gin-gonic/gin"

type MetricsControllerInterface interface {
	List(ctx *gin.Context)
	Show(ctx *gin.Context)
	Save(ctx *gin.Context)
	SaveBody(ctx *gin.Context)
	ShowBody(ctx *gin.Context)
	Ping(ctx *gin.Context)
	SaveMany(ctx *gin.Context)
}
