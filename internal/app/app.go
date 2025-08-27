package app

import (
	"log"
	"os"

	"github.com/augustoapg/censysKvStore/internal/api"
)

type Application struct {
	Logger    *log.Logger
	KVHandler *api.KVHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	return &Application{
		Logger:    logger,
		KVHandler: api.NewKVHandler(),
	}, nil
}
