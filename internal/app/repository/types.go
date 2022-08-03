package repository

import (
	qmgooptions "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ordered            = false
	insertIgnoreErrOpt = qmgooptions.InsertManyOptions{
		InsertHook: nil,
		InsertManyOptions: &options.InsertManyOptions{
			BypassDocumentValidation: nil,
			Ordered:                  &ordered,
		},
	}
)
