package testing

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/infinities-within/go-roomba"
	"github.com/infinities-within/go-roomba/sim"
)

var roombaSim *sim.RoombaSimulator
var mockRoombaClient *roomba.Roomba

func MakeTestRoomba() *roomba.Roomba {
	if mockRoombaClient == nil {
		var socket io.ReadWriter
		roombaSim, socket = sim.MakeRoombaSim()

		mockRoombaClient = &roomba.Roomba{S: socket, StreamPaused: make(chan bool, 1)}
	}
	return mockRoombaClient
}

func ClearTestRoomba() {
	mockRoombaClient = nil
	roombaSim.Stop()
	roombaSim = nil
}

func VerifyWritten(r *roomba.Roomba, expected []byte, t *testing.T) {
	// TODO : Rewrite xa4a's async read/write thing into something that's consistent
	time.Sleep(time.Millisecond * 100)
	actual := make([]byte, len(expected))
	roombaSim.ReadBytes.Read(actual)
	fmt.Println("Actual: ", actual)

	if len(actual) != len(expected) {
		t.Errorf("actual written length (%d) doesn't match expected (%d).",
			len(actual), len(expected))
	}
	for i, b := range expected {
		if b != actual[i] {
			t.Errorf("Expected output: % d, actual output: % d. Byte %d doesn't match",
				expected, actual, i)
		}
	}
}
