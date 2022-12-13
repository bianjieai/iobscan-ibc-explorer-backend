package rest

import (
	"fmt"
	"net/http"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/gin-gonic/gin"
)

type AddressController struct {
}

func (ctl *AddressController) BaseInfo(c *gin.Context) {
	chain := c.Param("chain")
	address := c.Param("address")

	if chain == "" || address == "" {
		c.JSON(http.StatusOK, response.FailBadRequest(fmt.Errorf("invalid parameters")))
		return
	}

	resp, err := addressService.BaseInfo(chain, address)
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}

func (ctl *AddressController) TokenList(c *gin.Context) {
	chain := c.Param("chain")
	address := c.Param("address")

	if chain == "" || address == "" {
		c.JSON(http.StatusOK, response.FailBadRequest(fmt.Errorf("invalid parameters")))
		return
	}

	resp, err := addressService.TokenList(chain, address)
	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(resp))
}
