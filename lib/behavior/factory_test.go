package behavior

import (
	"fmt"
	"testing"
)

func TestFactory_Register(t *testing.T) {
	fa := NewFactory()
	if err := fa.RegisterControlNode("test", NewSequenceNode); err != nil {
		t.Fatal(err)
	}
	fmt.Println(fa)
}
