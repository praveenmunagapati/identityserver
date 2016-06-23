package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateTime(t *testing.T) {
	now := time.Now()
	dt := DateTime(now)
	convertedNow := time.Time(dt)
	assert.Equal(t, now, convertedNow)
}
