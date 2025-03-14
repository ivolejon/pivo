package settings_test

import (
	"testing"

	"github.com/ivolejon/pivo/settings"
)

func TestGetEnv(t *testing.T) {
	_ = settings.Environment()
}
