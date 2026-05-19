//go:build windows

package notifications

import (
	"errors"

	"github.com/cryptdefender3232/phantom/server/configs"
	"github.com/nikoksr/notify"
)

func buildSyslog(_ *configs.SyslogConfig) (notify.Notifier, error) {
	return nil, errors.New("syslog notifications are not supported on windows")
}
