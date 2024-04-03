package sl

import "log/slog"

func Err(err error) slog.Attr {
	var errStr string = "no"

	if err != nil {
		errStr = err.Error()
	}

	return slog.String("error", errStr)
}
