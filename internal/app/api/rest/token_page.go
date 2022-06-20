package rest

import (
	"net/http"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/gin-gonic/gin"
)

type TokenController struct {
}

func (ctl *TokenController) List(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success(nil))
}
