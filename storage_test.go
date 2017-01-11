package mac

import "testing"

func TestStorage(t *testing.T) {
	var s storage
	t.Log(s.Resources())
	t.Log(s.CSS())
	t.Log(s.JS())
	t.Log(s.Default())
}
