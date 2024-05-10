package logs

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var Logger zerolog.Logger
var file *os.File

func NewLogger(pathFile string) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	file, err := os.OpenFile(
		pathFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}
	Logger = zerolog.New(file).With().Timestamp().Stack().Logger()
}
