package logging

import (
	"fmt"
	"testing"
)

func TestOS(t *testing.T) {
	t.Log(fmt.Sprintf("%c[1;%v;32m %s %c[0m", 0x1B, 43, "INFO", 0x1B))
}
