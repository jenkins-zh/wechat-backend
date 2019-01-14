package reply

import "testing"

func TestXML(t *testing.T) {
	data, err := makeTextResponseBody("", "", "")
	if err != nil {
		t.Errorf("xml error %v", err)
	}

	t.Errorf("%s", string(data))
}
