package digest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"gotest.tools/v3/assert"
)

func TestName(t *testing.T) {
	chalengeStr := `Digest  realm="asterisk",nonce="1726016274/381c86aa149ff6ec8eb72e704239926b",opaque="1622622838de7e4b",algorithm=md5,qop="auth"`
	challenge, err := ParseChallenge(chalengeStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(challenge)
}

func TestChallenge(t *testing.T) {
	tests := []struct {
		input     string
		challenge *Challenge
	}{
		{
			input: `Digest  realm="AXIS_ACCC8EB3494E", nonce="PNHWZB6nBQA=316099a140230c2db387fc75ee1c8ae838a750d8", stale=true, algorithm=MD5, qop="auth"`,
			challenge: &Challenge{
				Realm:     "AXIS_ACCC8EB3494E",
				Nonce:     "PNHWZB6nBQA=316099a140230c2db387fc75ee1c8ae838a750d8",
				Stale:     true,
				Algorithm: "MD5",
				QOP:       []string{"auth"},
			},
		},
		{
			input: `Digest  realm="AXIS_ACCC8EB3494E", nonce="PNHWZB6nBQA=316099a140230c2db387fc75ee1c8ae838a750d8", algorithm=MD5-sess, qop="auth"`,
			challenge: &Challenge{
				Realm:     "AXIS_ACCC8EB3494E",
				Nonce:     "PNHWZB6nBQA=316099a140230c2db387fc75ee1c8ae838a750d8",
				Algorithm: "MD5-sess",
				QOP:       []string{"auth"},
			},
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			c, err := ParseChallenge(tt.input)
			assert.NilError(t, err)
			assert.DeepEqual(t, tt.challenge, c)
			assert.DeepEqual(t, c.String(), tt.input)
		})
	}
}

func TestFindChallenge(t *testing.T) {
	bad1 := &Challenge{
		Realm:     "test",
		Nonce:     "kvjkdfjs",
		Algorithm: "MD5-sess",
		QOP:       []string{"auth"},
	}
	good := &Challenge{
		Realm:     "test",
		Nonce:     "jgdfsijdfisd",
		Algorithm: "MD5",
		QOP:       []string{"auth"},
	}
	headers := http.Header{}
	headers.Add("WWW-Authenticate", bad1.String())
	headers.Add("WWW-Authenticate", good.String())
	chal, err := FindChallenge(headers)
	assert.NilError(t, err)
	assert.DeepEqual(t, chal, good)
}

func TestFindChallenge_NotFound(t *testing.T) {
	_, err := FindChallenge(http.Header{})
	if !errors.Is(err, ErrNoChallenge) {
		t.Fatalf("not an expected error: %s, expected ErrNoChallenge", err.Error())
	}
}
