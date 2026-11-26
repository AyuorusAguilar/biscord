package token

import (
	"fmt"
	"testing"
	"time"
)

type test struct {
	id				int
	expectedToFail	bool
	expireIn		time.Duration
	waitTime		time.Duration
}

const secretDummy = "HolaBiscord"

var tests = []test{
	{id: 1, expectedToFail: true, expireIn: 2 * time.Second, waitTime: 3 * time.Second},
	{id: 1, expectedToFail: false, expireIn: 3 * time.Second, waitTime: 2 * time.Second},
	{id: 3, expectedToFail: false, expireIn: 3 * time.Second, waitTime: 2 * time.Second},
}

func TestJwts(ttool *testing.T) {
	for i, t := range tests {
		token, err := CreateJWT(t.id, secretDummy, t.expireIn)
		if err != nil {
			ttool.Fatalf("Error in test number %d:\n\t%v\n", i, err)
		}
		time.Sleep(t.waitTime)
		id, err := ValidateJWT(token, secretDummy)
		if err != nil {
			if t.expectedToFail { return }
			ttool.Fatalf("Error in test number %d:\n\t%v\n", i, err)
		}
		if id != int64(t.id) {
			ttool.Fatalf("Error in test number %d:\n\t%s\n", i, fmt.Sprintf("The returned id doesn't match! expecte:%d, got:%d", t.id, id))
		}
	}
}