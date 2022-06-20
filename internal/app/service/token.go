package service

type ITokenService interface {
	List()
}

type TokenService struct {
}

var _ ITokenService = new(TokenService)

func (svc *TokenService) List() {

}
