package service

type IChainService interface {
}

var _ IChainService = new(ChainService)

type ChainService struct {
}
