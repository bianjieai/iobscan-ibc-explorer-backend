package repository

type ITokenRepo interface {
}

var _ ITokenRepo = new(TokenRepo)

type TokenRepo struct {
}
