package task

import "github.com/bianjieai/iobscan-ibc-explorer-backend/internal/app/repository"

const (
	EveryMinute     = "0 */1 * * * ?"
	EverySecond     = "*/1 * * * * ?"
	EveryFiveSecond = "*/5 * * * * ?"
	EveryTenSecond  = "*/10 * * * * ?"
	OneHour         = "0 0 */1 * * ?"
	TwelveHour      = "0 0 */12 * * ?"
	ThreeMinute     = "0 */3 * * * ?"
	FiveMinute      = "0 */5 * * * ?"
	TwentyMinute    = "0 */20 * * * ?"
)

var (
	tokenRepo    repository.ITokenRepo       = new(repository.TokenRepo)
	chainCfgRepo repository.IChainConfigRepo = new(repository.ChainConfigRepo)
	chainRepo    repository.IChainRepo       = new(repository.IbcChainRepo)
	ibcTxRepo    repository.IExIbcTxRepo     = new(repository.ExIbcTxRepo)
	relayerRepo  repository.IRelayerRepo     = new(repository.IbcRelayerRepo)
)
