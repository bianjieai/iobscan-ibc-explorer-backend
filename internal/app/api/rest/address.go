package rest

import (
	"fmt"
	"net/http"

	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/api/response"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/errors"
	"github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/model/vo"
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

func (ctl *AddressController) TxsList(c *gin.Context) {
	chain := c.Param("chain")
	address := c.Param("address")

	if chain == "" || address == "" {
		c.JSON(http.StatusOK, response.FailBadRequest(fmt.Errorf("invalid parameters")))
		return
	}

	var req vo.AddressTxsListReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, response.FailBadRequest(err))
		return
	}

	var res interface{}
	var err errors.Error
	if !req.UseCount {
		res, err = addressService.TxsList(chain, address, &req)
	} else {
		res, err = addressService.TxsCount(chain, address)
	}

	if err != nil {
		c.JSON(http.StatusOK, response.FailError(err))
		return
	}
	c.JSON(http.StatusOK, response.Success(res))
}

func (ctl *AddressController) TxsExport(c *gin.Context) {
	chain := c.Param("chain")
	address := c.Param("address")

	if chain == "" || address == "" {
		c.JSON(http.StatusBadRequest, response.FailBadRequest(fmt.Errorf("invalid parameters")))
		return
	}

	filename, data, err := addressService.TxsExport(chain, address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.FailError(err))
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", filename))
	c.Data(http.StatusOK, "text/csv", data)
}
