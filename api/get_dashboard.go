package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getReport(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "Report")
}
