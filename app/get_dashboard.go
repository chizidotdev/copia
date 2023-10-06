package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) getReport(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "Report")
}
