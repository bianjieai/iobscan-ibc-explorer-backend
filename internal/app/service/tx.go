package service

type ITxService interface {
}

var _ ITxService = new(TxService)

type TxService struct {
}
