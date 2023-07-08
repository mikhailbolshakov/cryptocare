package http

import (
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	kitTestSuite "github.com/mikhailbolshakov/cryptocare/src/kit/test/suite"
	"github.com/stretchr/testify/suite"
	"testing"
)

var logger = log.Init(&log.Config{Level: log.InfoLevel})
var logf = func() log.CLogger {
	return log.L(logger)
}

type sortConvertTestSuite struct {
	kitTestSuite.Suite
}

func (s *sortConvertTestSuite) SetupSuite() {
	s.Suite.Init(logf)
}

func TestTagSuite(t *testing.T) {
	suite.Run(t, new(sortConvertTestSuite))
}

func (s *sortConvertTestSuite) Test_ParseSortBy() {
	tests := []struct {
		name       string
		sortString string
		want       []*SortRequest
		wantErr    bool
	}{
		{
			name:       "Empty string",
			sortString: "",
			want:       nil,
		},
		{
			name:       "real example",
			sortString: "reportedAt desc",
			want: []*SortRequest{
				{
					Field: "reportedAt",
					Asc:   false,
				},
			},
		},
		{
			name:       "All ok (without missings)",
			sortString: "field1,field2 desc",
			want: []*SortRequest{
				{
					Field: "field1",
					Asc:   true,
				},
				{
					Field: "field2",
					Asc:   false,
				},
			},
		},
		{
			name:       "All ok (with missings)",
			sortString: "field1 asc first,field2 desc last,field3 asc",
			want: []*SortRequest{
				{
					Field:   "field1",
					Asc:     true,
					Missing: SortRequestMissingFirst,
				},
				{
					Field:   "field2",
					Asc:     false,
					Missing: SortRequestMissingLast,
				},
				{
					Field:   "field3",
					Asc:     true,
					Missing: "",
				},
			},
		},
		{
			name:       "Whitespaces",
			sortString: " field1    asc  , field2 desc  ",
			wantErr:    true,
		},
		{
			name:       "1 field",
			sortString: "field1 asc",
			want: []*SortRequest{
				{
					Field: "field1",
					Asc:   true,
				},
			},
		},
		{
			name:       "1 field only name",
			sortString: "field1",
			want: []*SortRequest{
				{
					Field: "field1",
					Asc:   true,
				},
			},
		},
		{
			name:       "Illegal sort mode",
			sortString: "field1 asc,field2 illegal_mode",
			wantErr:    true,
		},
		{
			name:       "Illegal missing mode",
			sortString: "field1 asc,field2 desc illegal_mode",
			wantErr:    true,
		},
		{
			name:       "Illegal syntax 1",
			sortString: "field1 asc,field2=desc",
			wantErr:    true,
		},
		{
			name:       "Illegal syntax 2",
			sortString: "field1 asc,desc=field2",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		res, err := ParseSortBy(s.Ctx, tt.sortString)
		s.Equal(tt.want, res)
		s.Equal(tt.wantErr, err != nil)
	}
}
