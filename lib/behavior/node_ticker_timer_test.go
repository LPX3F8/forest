package behavior

import (
	"fmt"
	"testing"
	"time"
)

func TestPeriod(t *testing.T) {
	p := new(Period)
	p.Start()
	time.Sleep(time.Second)
	p.Stop()
	fmt.Println(p)
}
