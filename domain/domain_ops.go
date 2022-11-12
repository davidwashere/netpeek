package domain

import (
	"encoding/json"
	"fmt"
)

func (r Result) Bytes() []byte {
	dataB, _ := json.Marshal(r)
	return dataB
}

func (r Result) String() string {
	return string(r.Bytes())
}

func (r Result) PrettyString() string {
	dataB, _ := json.MarshalIndent(r, "", "  ")
	return string(dataB)
}

func (t *NetTuple) Key() string {
	return fmt.Sprintf("%s:%s->%s:%s", t.SrcIP, t.SrcPort, t.DestIP, t.DestPort)
}
