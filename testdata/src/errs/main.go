package a

import (
	"fmt"
	"log/slog"
)

func main(tk, secret, bearer string) {
	slog.Info("Hello")

	slog.Info("привет")

	slog.Info("bad!")

	slog.Info("token expired")

	slog.Info("token=" + tk)

	slog.Info(fmt.Sprintf("token=%s", tk))

	slog.Info(fmt.Sprintf("secret=%+v", secret))

	slog.Info(fmt.Sprintf("Authorization: Bearer %s", bearer))

	slog.Info("token=")

	slog.Info(fmt.Sprintf("token=%s"))
}
