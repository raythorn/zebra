package log

import (
	"fmt"
	"github.com/raythorn/zebra/log"
	"testing"
)

func TestLog(t *testing.T) {

	fmt.Println("Start testing log...")

	log.Debug("Test Debug")
	log.Info("Test Info")
	log.Warning("Test Warning")
}
