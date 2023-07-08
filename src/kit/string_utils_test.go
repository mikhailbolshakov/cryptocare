package kit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Strings_Distinct(t *testing.T) {
	s := Strings{"aaa", "bbb", "aaa"}
	assert.Equal(t, s.Distinct(), Strings{"aaa", "bbb"})
	s = Strings{"aaa", "bbb"}
	assert.Equal(t, s.Distinct(), Strings{"aaa", "bbb"})
}

func Test_Strings_Contains(t *testing.T) {
	s := Strings{"aaa", "bbb", "aaa"}
	assert.True(t, s.Contains("aaa"))
	assert.True(t, s.Contains("bbb"))
	assert.False(t, s.Contains("ccc"))
	assert.False(t, s.Contains(""))
	assert.False(t, s.Contains("aa"))
}

func Test_Strings_Intersect(t *testing.T) {
	for _, s := range []struct {
		S1 Strings
		S2 Strings
		R  Strings
	}{
		{
			S1: nil,
			S2: nil,
			R:  nil,
		},
		{
			S1: Strings{},
			S2: nil,
			R:  nil,
		},
		{
			S1: Strings{},
			S2: Strings{},
			R:  nil,
		},
		{
			S1: Strings{"aa"},
			S2: Strings{},
			R:  nil,
		},
		{
			S1: Strings{"aa"},
			S2: Strings{"bb"},
			R:  nil,
		},
		{
			S1: Strings{"aa", "bb"},
			S2: Strings{"bb", "cc"},
			R:  Strings{"bb"},
		},
		{
			S1: Strings{"aa", "bb", "dd"},
			S2: Strings{"bb", "cc", "dd"},
			R:  Strings{"bb", "dd"},
		},
	} {
		assert.Equal(t, s.R, s.S1.Intersect(s.S2))
	}
}

func Test_Strings_Equal(t *testing.T) {
	for _, s := range []struct {
		S1 Strings
		S2 Strings
		R  bool
	}{
		{
			S1: nil,
			S2: nil,
			R:  true,
		},
		{
			S1: Strings{},
			S2: nil,
			R:  true,
		},
		{
			S1: Strings{},
			S2: Strings{},
			R:  true,
		},
		{
			S1: Strings{"aa"},
			S2: Strings{},
			R:  false,
		},
		{
			S1: Strings{"aa"},
			S2: Strings{"bb"},
			R:  false,
		},
		{
			S1: Strings{"aa", "bb"},
			S2: Strings{"bb", "cc"},
			R:  false,
		},
		{
			S1: Strings{"aa", "bb"},
			S2: Strings{"aa", "bb", "cc"},
			R:  false,
		},
		{
			S1: Strings{"aa", "bb", "cc"},
			S2: Strings{"aa", "bb", "cc"},
			R:  true,
		},
	} {
		assert.Equal(t, s.R, s.S1.Equal(s.S2))
	}
}

func Test_Strings_Subset(t *testing.T) {
	for _, s := range []struct {
		S1 Strings
		S2 Strings
		R  bool
	}{
		{
			S1: nil,
			S2: nil,
			R:  true,
		},
		{
			S1: Strings{},
			S2: nil,
			R:  true,
		},
		{
			S1: Strings{},
			S2: Strings{},
			R:  true,
		},
		{
			S1: Strings{"aa"},
			S2: Strings{},
			R:  false,
		},
		{
			S1: Strings{},
			S2: Strings{"aa"},
			R:  true,
		},
		{
			S1: Strings{"aa"},
			S2: Strings{"bb"},
			R:  false,
		},
		{
			S1: Strings{"aa", "bb"},
			S2: Strings{"bb", "cc"},
			R:  false,
		},
		{
			S1: Strings{"aa", "bb", "cc"},
			S2: Strings{"aa", "bb"},
			R:  false,
		},
		{
			S1: Strings{"aa", "bb", "cc"},
			S2: Strings{"aa", "bb", "cc"},
			R:  true,
		},
		{
			S1: Strings{"aa", "bb"},
			S2: Strings{"aa", "bb", "cc"},
			R:  true,
		},
	} {
		assert.Equal(t, s.R, s.S1.Subset(s.S2))
	}
}

func Test_Strings_Sanitize(t *testing.T) {
	for _, s := range []struct {
		S Strings
		R Strings
	}{
		{
			S: nil,
			R: nil,
		},
		{
			S: Strings{},
			R: Strings{},
		},
		{
			S: Strings{"aa"},
			R: Strings{"aa"},
		},
		{
			S: Strings{" Aa-.^&(^&	___"},
			R: Strings{"aa"},
		},
		{
			S: Strings{" Aa-.^&(^&	___ ", "  242   %%%%---####   "},
			R: Strings{"aa", "242"},
		},
	} {
		assert.Equal(t, s.R, s.S.Sanitize())
	}
}

func Test_StrToInt64(t *testing.T) {
	for _, s := range []struct {
		In  string
		Out int64
		Err bool
	}{
		{
			In:  "",
			Out: 0,
			Err: true,
		},
		{
			In:  "qqq",
			Out: 0,
			Err: true,
		},
		{
			In:  "0.23123",
			Out: 0,
			Err: true,
		},
		{
			In:  "-1",
			Out: -1,
			Err: false,
		},
		{
			In:  "1576663112362381",
			Out: 1576663112362381,
			Err: false,
		},
	} {
		out, err := StrToInt64(s.In)
		if s.Err {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, s.Out, out)
		}
	}
}

func Test_RemoveNonAlfaDigital(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "empty string",
			in:   "",
			out:  "",
		},
		{
			name: "complex case",
			in: "  A++B%%C///--	 %:*%*abc \t@#$%123   &&& АБВ *)__^^ абв",
			out: "ABCabc123АБВабв",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, RemoveNonAlfaDigital(tt.in))
		})
	}
}
