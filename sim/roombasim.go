/*
Package sim defines a limited OpenInterface simulator type that's mainly used for testing.

Simulator can be created using MakeRoombaSim() function, which returns a
simulator instance and a ReadWriter, suitable for passing to go-roomba client.
*/
package sim

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/infinities-within/go-roomba"
	"github.com/infinities-within/go-roomba/constants"
)

// Roomba simulator instance. Should be constructed with MakeRoombaSim()
// function.
type RoombaSimulator struct {
	rw           io.ReadWriter
	writeQ       chan []byte
	WrittenBytes bytes.Buffer // Logs all the bytes written by the simulator to its Writer.
	ReadBytes    bytes.Buffer // Logs all the bytes read by the simulator from its Reader.

	RequestedVelocity []byte
	RequestedRadius   []byte
}

// MockSensorValues contains mapping of sensor codes to sensor values returned
// by a RoombaSimulator object on sensor requests.
var MockSensorValues = map[constants.SensorCode][]byte{
	constants.SENSOR_BUMP_WHEELS_DROPS:       []byte{3},
	constants.SENSOR_VIRTUAL_WALL:            []byte{5},
	constants.SENSOR_CLIFF_RIGHT:             []byte{42},
	constants.SENSOR_TEMPERATURE:             []byte{25},
	constants.SENSOR_OI_MODE:                 []byte{2},
	constants.SENSOR_SONG_NUMBER:             []byte{1},
	constants.SENSOR_DISTANCE:                []byte{10, 20},
	constants.SENSOR_WALL:                    []byte{35},
	constants.SENSOR_BATTERY_CHARGE:          roomba.Pack([]interface{}{uint16(1000)}),
	constants.SENSOR_BATTERY_CAPACITY:        roomba.Pack([]interface{}{uint16(1500)}),
	constants.SENSOR_CURRENT:                 roomba.Pack([]interface{}{int16(-747)}),
	constants.SENSOR_CLIFF_FRONT_LEFT_SIGNAL: roomba.Pack([]interface{}{uint8(2), uint8(25)}),
}

func (sim *RoombaSimulator) serve() {
	// Write bytes from channel asynchronously.
	go func() {
		for {
			bs := <-sim.writeQ
			if len(bs) == 0 {
				break
			}
			sim.rw.Write(bs)
		}
	}()

	for {
		sim.executeCMD()
	}
}

func (sim *RoombaSimulator) Stop() {
	sim.writeQ <- []byte{}
}

func (sim *RoombaSimulator) executeCMD() error {
	cmdBuf := sim.read(1)
	if len(cmdBuf) != 1 {
		return fmt.Errorf("failed reading opcode")
	}
	switch constants.OpCode(cmdBuf[0]) {
	case constants.Sensors:
		packetId := constants.SensorCode(sim.read(1)[0])
		value, ok := MockSensorValues[packetId]
		if !ok {
			if packetId == constants.SENSOR_REQUESTED_RADIUS {
				value = sim.RequestedRadius
			} else if packetId == constants.SENSOR_REQUESTED_VELOCITY {
				value = sim.RequestedVelocity
			} else {
				log.Printf("no mock value for sensor packet id %d", packetId)
			}
		}
		log.Printf("sensor %d value: %v", packetId, value)
		sim.write(value)
	case constants.QueryList:
		nPackets := sim.read(1)[0]
		for i := 0; i < int(nPackets); i++ {
			packetId := constants.SensorCode(sim.read(1)[0])
			value, ok := MockSensorValues[packetId]
			if !ok {
				if packetId == constants.SENSOR_REQUESTED_RADIUS {
					value = sim.RequestedRadius
				} else if packetId == constants.SENSOR_REQUESTED_VELOCITY {
					value = sim.RequestedVelocity
				} else {
					log.Printf("no mock value for sensor packet id %d", packetId)
				}
			}
			log.Printf("sensor %d value: %v", packetId, value)
			sim.write(value)
		}
	case constants.SensorStream:
		nBytes := sim.read(1)[0]
		packetIds := make([]constants.SensorCode, nBytes)
		for i := byte(0); i < nBytes; i++ {
			packetIds[i] = constants.SensorCode(sim.read(1)[0])
		}
		// Contains just packet ids and values, no headers.
		sensorValues := bytes.Buffer{}
		for i := byte(0); i < nBytes; i++ {
			mockValue, ok := MockSensorValues[packetIds[i]]
			if !ok {
				log.Printf("no mock value for streaming packet id: %d", packetIds[i])
				mockValue = make([]byte, constants.SENSOR_PACKET_LENGTH[packetIds[i]])
			} else {
				log.Printf("sensor %d value: %v", packetIds[i], mockValue)
			}
			sensorValues.WriteByte(byte(packetIds[i]))
			sensorValues.Write(mockValue)
		}

		output := bytes.Buffer{}
		// Header.
		output.WriteByte(19)
		// Data length.
		messageLen := byte(sensorValues.Len())
		log.Printf("message length: %d", messageLen)
		output.WriteByte(messageLen)
		output.Write(sensorValues.Bytes())
		checksum := byte(0)
		for _, b := range output.Bytes()[1:] {
			checksum -= b
		}
		output.WriteByte(checksum)
		log.Printf("checksum: %d", checksum)

		sim.write(output.Bytes())
	case constants.Start:
		log.Printf("switched to passive mode")
	case constants.Safe:
		log.Printf("switched to safe mode")
	case constants.PauseResumeStream:
		if sim.read(1)[0] == byte(0) {
			log.Printf("stream paused")
		} else {
			log.Printf("stream resumed")
		}
	case constants.DriveDirect:
		data := sim.read(4)
		var rightVelocity, leftVelocity int16
		_ = binary.Read(bytes.NewReader(data[:2]), binary.BigEndian, &rightVelocity)
		_ = binary.Read(bytes.NewReader(data[2:4]), binary.BigEndian, &leftVelocity)
		log.Printf("DirectDrive: %d, %d (%v)", rightVelocity, leftVelocity, data)
	case constants.Drive:
		sim.RequestedVelocity = sim.read(2)
		sim.RequestedRadius = sim.read(2)
		log.Printf("Drive: %d, %d", sim.RequestedVelocity, sim.RequestedRadius)
	default:
		log.Printf("unknown opcode: %d", cmdBuf[0])
	}

	return nil
}

// Reads given number of bytes from the Reader sim.rw.
func (sim *RoombaSimulator) read(n int) []byte {
	buf := make([]byte, n)
	nRead, err := sim.rw.Read(buf)
	if n != nRead {
		if err != nil {
			log.Printf("error reading in RoombaSimulator: %v", err)
		}
		//log.Printf("read %d bytes when expected %d", nRead, n)
		return []byte{}
	}
	log.Printf("roomba reads: %v", buf)
	sim.ReadBytes.Write(buf)
	return buf
}

// Writes bytes to the Writer w asynchronously.
func (sim *RoombaSimulator) write(b []byte) {
	log.Printf("roomba says: %v", b)
	sim.writeQ <- b
}

// Helper for merging reader and writer into a ReadWriter.
type readWriter struct {
	io.Reader
	io.Writer
}

func MakeRoombaSim() (*RoombaSimulator, *readWriter) {
	// Input: driver writes, simulator reads.
	inp_r, inp_w := io.Pipe()

	// Ouput: simulator writes, driver reads.
	out_r, out_w := io.Pipe()

	readBytes := &bytes.Buffer{}
	writtenBytes := &bytes.Buffer{}

	sim := &RoombaSimulator{
		rw: &readWriter{
			// Log all read bytes to ReadBytes.
			io.TeeReader(inp_r, readBytes),
			// Log all written bytes to writtenBytes.
			io.MultiWriter(out_w, writtenBytes),
		},
		writeQ:    make(chan []byte, 15),
		ReadBytes: *readBytes,

		RequestedRadius:   []byte{0, 0},
		RequestedVelocity: []byte{0, 0},
	}
	go sim.serve()

	rw := &readWriter{out_r, inp_w}

	return sim, rw
}
