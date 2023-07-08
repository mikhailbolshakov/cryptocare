package kit

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_TimeParse(t *testing.T) {
	tm := &Time{}
	err := tm.Parse("10:00")
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, tm.Time.String())
	assert.Equal(t, 10, tm.Time.Hour())
	assert.Equal(t, 0, tm.Time.Minute())
}

func Test_TimeString(t *testing.T) {
	tm := &Time{}
	err := tm.Parse("10:00")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "10:00", tm.String())
}

func Test_Date(t *testing.T) {
	tests := []struct {
		name   string
		input  time.Time
		output time.Time
	}{
		{
			name:   "Time to date",
			input:  time.Date(2022, 1, 1, 12, 23, 34, 122, time.UTC),
			output: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:   "Date to date",
			input:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			output: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.output, Date(tt.input))
		})
	}
}
