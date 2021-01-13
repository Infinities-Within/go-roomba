// Package roomba iRobot roomba Open Interface.
//
// The Roomba OI has four operating modes: Off, Passive, Safe, and Full.
// When roomba starts the OI is in “off” mode. When it is off, the OI listens
// for an OI Start command. Once it receives the Start command, you can enter
// into any one of the four operating modes by sending a mode command to the OI.
//
// Passive mode: entered upon sending one of the cleaning commands. You can
// only read sensor data in the passive mode and can't change the actuators
// state.
//
// Safe mode: gives full control of Roomba, except for safety restrictions:
//   * Cliff detection when moving forward.
//	 * Detection of wheel drop.
// 	 * Charger plugged in and powered.
// When any of the events ocurs, Roomba switches to passive mode.
//
// Full mode: gives full control over Romoba, disabling the safety
// restrictions.

package roomba

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/infinities-within/go-roomba/constants"
)

func toByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

// MakeRoomba initializes a new Roomba structure and sets up a serial port.
// By default, Roomba communicates at 57600 baud.
func MakeRoomba(portName string) (*Roomba, error) {
	roomba := &Roomba{PortName: portName, StreamPaused: make(chan bool, 1)}
	baud := uint(57600)
	err := roomba.Open(baud)
	return roomba, err
}

// Start command starts the OI. You must always send the Start command before
// sending any other commands to the OI.
// Note: Use the Start command (128) to change the mode to Passive.
func (roomba *Roomba) Start() error {
	return roomba.WriteByte(constants.Start)
}

// TODO: Baud command.

// Passive switches Roomba to passive mode by sending the Start command.
func (roomba *Roomba) Passive() error {
	return roomba.Start()
}

// This command puts the OI into Safe mode, enabling user control of Roomba.
// It turns off all LEDs.
func (roomba *Roomba) Safe() error {
	return roomba.WriteByte(constants.Safe)
}

// Full command gives you complete control over Roomba by putting the OI into
// Full mode, and turning off the cliff, wheel-drop and internal charger safety
// features.
func (roomba *Roomba) Full() error {
	return roomba.WriteByte(constants.Full)
}

// Control command's effect and usage are identical to the Safe command.
func (roomba *Roomba) Control() error {
	roomba.Passive()
	return roomba.WriteByte(constants.Control) // ?
}

// Clean command starts the default cleaning mode.
func (roomba *Roomba) Clean() error {
	return roomba.WriteByte(constants.Cover)
}

// TODO: Max command.

// Spot command starts the Spot cleaning mode.
func (roomba *Roomba) Spot() error {
	return roomba.WriteByte(constants.Spot)
}

// SeekDock command sends Roomba to the dock.
func (roomba *Roomba) SeekDock() error {
	return roomba.WriteByte(constants.Dock)
}

// Drive command controls Roomba’s drive wheels. It takes two 16-bit signed
// values. The first one specifies the average velocity of the drive wheels in
// millimeters per second (mm/s).  The next one specifies the radius in
// millimeters at which Roomba will turn. The longer radii make Roomba drive
// straighter, while the shorter radii make Roomba turn more. The radius is
// measured from the center of the turning circle to the center of Roomba. A
// Drive command with a positive velocity and a positive radius makes Roomba
// drive forward while turning toward the left. A negative radius makes Roomba
// turn toward the right. Special cases for the radius make Roomba turn in place
// or drive straight. A negative velocity makes Roomba drive backward. Velocity
// is in range (-500 – 500 mm/s), radius (-2000 – 2000 mm). Special cases:
// straight = 32768 or 32767 = hex 8000 or 7FFF, turn in place clockwise = -1,
// turn in place counter-clockwise = 1
func (roomba *Roomba) Drive(velocity, radius int16) error {
	if !(-500 <= velocity && velocity <= 500) {
		return fmt.Errorf("invalid velocity: %d", velocity)
	}
	if !(-2000 <= radius && radius <= 2000) {
		return fmt.Errorf("invalid readius: %d", radius)
	}
	return roomba.Write(constants.Drive, Pack([]interface{}{velocity, radius}))
}

// Stop commands is equivalent to Drive(0, 0).
func (roomba *Roomba) Stop() error {
	return roomba.Drive(0, 0)
}

// DirectDrive command lets you control the forward and backward motion of
// Roomba’s drive wheels independently. It takes two 16-bit signed values.
// The first specifies the velocity of the right wheel in millimeters per second
// (mm/s), The next one specifies the velocity of the left wheel A positive
// velocity makes that wheel drive forward, while a negative velocity makes it
// drive backward. Right wheel velocity (-500 – 500 mm/s). Left wheel velocity
// (-500 – 500 mm/s).
func (roomba *Roomba) DirectDrive(right, left int16) error {
	if !(-500 <= right && right <= 500) ||
		!(-500 <= left && left <= 500) {
		return fmt.Errorf("invalid velocity. one of %d or %d", right, left)
	}
	return roomba.Write(constants.DriveDirect, Pack([]interface{}{right, left}))
}

// TODO: Drive PWM, Motors, PWM Motors commands.

// LEDs command controls the LEDs common to all models of Roomba 500. The
// Clean/Power LED is specified by two data bytes: one for the color and the
// other for the intensity. Color: 0 = green, 255 = red. Intermediate values are
// intermediate colors (orange, yellow, etc). Intensitiy: 0 = off, 255 = full
// intensity. Intermediate values are intermediate intensities.
func (roomba *Roomba) LEDs(advance, play bool, powerColor, powerIntensity byte) error {
	var ledBits byte

	if advance {
		ledBits += 8
	}
	if play {
		ledBits += 2
	}

	return roomba.Write(constants.LEDs, Pack([]interface{}{
		ledBits, powerColor, powerIntensity}))
}

// Sensors command requests the OI to send a packet of sensor data bytes. There
// are 58 different sensor data packets. Each provides a value of a specific
// sensor or group of sensors.
func (roomba *Roomba) Sensors(packetId constants.SensorCode) ([]byte, error) {
	bytesToRead, ok := constants.SENSOR_PACKET_LENGTH[packetId]
	if !ok {
		return []byte{}, fmt.Errorf("unknown packet id requested: %d", packetId)
	}

	roomba.Write(constants.Sensors, []byte{byte(packetId)})
	var err error
	var n int
	result := make([]byte, bytesToRead)
	for byte(n) < bytesToRead {
		resultView := result[n:]
		bytesToRead -= byte(n)
		n, err = roomba.Read(resultView)
		if err != nil {
			log.Printf("error %v", err)
			return result, fmt.Errorf("failed reading sensors data for packet id %d: %s", packetId, err)
		}
	}
	return result, nil
}

// QueryList command lets you ask for a list of sensor packets. The result is
// returned once, as in the Sensors command. The robot returns the packets in
/// the order you specify.
func (roomba *Roomba) QueryList(packetIds []constants.SensorCode) ([][]byte, error) {
	for _, packetId := range packetIds {
		_, ok := constants.SENSOR_PACKET_LENGTH[packetId]
		if !ok {
			return [][]byte{}, fmt.Errorf("unknown packet id requested: %d", packetId)
		}
	}

	b := new(bytes.Buffer)
	b.WriteByte(byte(len(packetIds)))
	for _, id := range packetIds {
		b.WriteByte(byte(id))
	}
	roomba.Write(constants.QueryList, b.Bytes())

	var err error
	var n int
	result := make([][]byte, len(packetIds))
	for i, packetId := range packetIds {
		bytesToRead := constants.SENSOR_PACKET_LENGTH[packetId]
		result[i] = make([]byte, bytesToRead)
		err, n = nil, 0
		for byte(n) < bytesToRead {
			resultView := result[i][n:]
			bytesToRead -= byte(n)
			n, err = roomba.Read(resultView)
			if err != nil {
				return result, fmt.Errorf("failed reading sensors data for packet id %d: %s", packetId, err)
			}
		}
	}
	return result, nil
}

// PauseStream command lets you stop steam without clearing the list of
// requested packets.
func (roomba *Roomba) PauseStream() {
	roomba.StreamPaused <- true
}

func (roomba *Roomba) ReadStream(packetIds []constants.SensorCode, out chan<- [][]byte) {
	var dataLength byte
	for _, packetId := range packetIds {
		packetLength, ok := constants.SENSOR_PACKET_LENGTH[packetId]
		if !ok {
			log.Printf("unknown packet id requested: %d", packetId)
			return
		}
		dataLength += packetLength
	}

	// Input buffer. 3 is for 19, N-bytes and checksum.
	buf := make([]byte, dataLength+byte(len(packetIds))+3)

	for {
	Loop:
		select {
		case <-roomba.StreamPaused:
			// Pause stream.
			roomba.Write(constants.PauseResumeStream, []byte{0})
			close(out)
			return
		default:
			// Read single stream frame.
			bytesRead := 0
			for bytesRead < len(buf) {
				n, err := roomba.S.Read(buf[bytesRead:])
				if n != 0 {
					bytesRead += n
				}
				if err != nil {
					if err == io.EOF {
						return
					}
					goto Loop
				}
			}
			// Process frame.
			bufR := bytes.NewReader(buf)
			if b, err := bufR.ReadByte(); err != nil || b != 19 {
				log.Fatalf("stream data doesn't start with header 19")
				return
			}
			if b, err := bufR.ReadByte(); err != nil || b != byte(len(buf)-3) {
				log.Fatalf("invalid N-bytes: %d, expected %d.", buf[1],
					len(buf)-3)
			}

			result := make([][]byte, len(packetIds))

			i := 0
			// Used for verifying checksum.
			sum := byte(len(buf) - 3) // N-bytes
			packetId, err := bufR.ReadByte()
			for ; err == nil; packetId, err = bufR.ReadByte() {
				sum += packetId
				bytesToRead := int(constants.SENSOR_PACKET_LENGTH[constants.SensorCode(packetId)])
				bytesRead := 0
				result[i] = make([]byte, bytesToRead)

				for bytesToRead > 0 {
					n, err := bufR.Read(result[i][bytesRead:])
					bytesRead += n
					bytesToRead -= n
					if err != nil {
						log.Fatalf("error reading packet data")
					}
				}
				for _, b := range result[i] {
					sum += b
				}
				i += 1
				if bufR.Len() == 1 {
					break
				}
			}

			expectedChecksum, err := bufR.ReadByte()
			if err != nil {
				log.Fatalf("missing checksum")
			}
			sum += expectedChecksum
			if sum != 0 {
				log.Fatalf("computed checksum didn't match: %d", sum)
			}
			out <- result
		}
	}
}

// Stream command starts a stream of data packets. The list of packets
// requested is sent every 15 ms, which is the rate Roomba uses to update data.
// This method of requesting sensor data is best if you are controlling Roomba
// over a wireless network (which has poor real-time characteristics) with
// software running on a desktop computer.
func (roomba *Roomba) Stream(packetIds []constants.SensorCode) (<-chan [][]byte, error) {
	b := new(bytes.Buffer)
	b.WriteByte(byte(len(packetIds)))
	for _, pid := range packetIds {
		b.WriteByte(byte(pid))
	}
	err := roomba.Write(constants.SensorStream, b.Bytes())
	if err != nil {
		return nil, err
	}

	out := make(chan [][]byte)
	go roomba.ReadStream(packetIds, out)
	return out, nil
}
