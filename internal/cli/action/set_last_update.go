package action

import (
	"github.com/spf13/viper"
	"time"
)

func SetLastUpdate() {
	viper.Set("lastupdate", time.Now().Unix())
}
