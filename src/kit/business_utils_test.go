package kit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IsEmailValid(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  bool
	}{
		{name: "valid", in: "test@test.com", out: true},
		{name: "with dot in domain", in: "test@test.test.com", out: true},
		{name: "with dot", in: "test.test@test.com", out: true},
		{name: "with plus and minus signs", in: "test.test+test-test@test.com", out: true},
		{name: "with quotes around username", in: "\"test\"@test.com", out: true},
		{name: "with quotes and @ inside username", in: "\"test.@.test\"@test.com", out: true},
		{name: "with allowed special symbols", in: "#!$%&'*+-/=?^_`{}|~@test.com", out: true},
		{name: "cyrillic", in: "тест@тест.рф", out: true},
		{name: "cyrillic username with latin domain", in: "тест@test.com", out: true},
		{name: "latin username with cyrillic domain", in: "test@тест.рф", out: true},
		{name: "too long username", in: "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklm@test.com", out: true},
		{name: "not valid tld", in: "3@test.abcdefghij", out: true},
		{name: "without domain", in: "test@test", out: true},

		{name: "double @", in: "te@st@test.com", out: false},
		{name: "with open brace", in: "tes(t@test.com", out: false},
		{name: "with close brace", in: "tes)t@test.com", out: false},
		{name: "with <", in: "tes<t@test.com", out: false},
		{name: "with >", in: "tes>t@test.com", out: false},
		{name: "with comma", in: "tes,t@test.com", out: false},
		{name: "with colon", in: "tes:t@test.com", out: false},
		{name: "with semicolon", in: "tes;t@test.com", out: false},
		{name: "empty email", in: "", out: false},
		{name: "dot at the end", in: "test@test.com.", out: false},
		{name: "dot at the end of username", in: "test.@test.com", out: false},
		{name: "dot at the beginning of username", in: ".test@test.com", out: false},
		{name: "too long hostname", in: "1@abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijfg.com", out: false},
		{name: "too long domain", in: "2@test.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijfghabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijfghabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijfghabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzab", out: false},
		{name: "without @", in: "testtest.com", out: false},
		{name: "without username", in: "@test.com", out: false},
		{name: "with multiple @", in: "A@b@c@test.com", out: false},
		{name: "with quotes inside username", in: "just\"not\"right@test.com", out: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, IsEmailValid(tt.in))
		})
	}
}
