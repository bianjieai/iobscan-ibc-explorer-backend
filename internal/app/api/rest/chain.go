package rest

import (
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChainController struct {
}

func (ctl *ChainController) List(c *gin.Context) {

	// todo
	c.JSON(http.StatusOK, response.Success(nil))
}
