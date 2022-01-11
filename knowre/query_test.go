package knowre

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestMakeQuery(t *testing.T) {
	t1 := time.Now()
	match := "/ecs/krdky-stable"
	i := 0
	dslQuery, query, err := MakeQuery(t1, match, 1, i)
	_ = query
	assert.Nil(t, err)
	log.Println(dslQuery)
}
