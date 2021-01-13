// Package constants defines values for OpenInterface op codes, sensor codes and sensor packet lengths among others.
package constants

type OpCode byte

const (
    Start = OpCode(iota + 128)
    Baud
    Control
    Safe
    Full
    _ // 133 is unused
    Spot
    Cover
    Demo
    Drive
    LowSideDrivers
    LEDs
    Song
    Play
    Sensors
    Dock
    PWMLowSideDrivers
    DriveDirect
    _ // 146 unused
    DigitalOutputs
    SensorStream
    QueryList
    PauseResumeStream
    SendIR
    Script
    PlayScript
    ShowScript
    WaitTime
    WaitDistance
    WaitAngle
    WaitEvent
)

type SensorCode byte

// SENSOR_* constants define the packet IDs for declared sensor packets.
const (
    // The state of the bumper (0 = no bump, 1 = bump) and wheel drop sensors
    // (0 = wheel raised, 1 = wheel dropped) are sent as individual bits.
    SENSOR_BUMP_WHEELS_DROPS = SensorCode(iota + 7)

    // The state of the wall sensor is sent as a 1 bit value (0 = no wall,
    // 1 = wall seen).
    SENSOR_WALL

    // The state of the cliff sensor on the left side of Roomba is sent as a 1
    // bit value (0 = no cliff, 1 = cliff).
    SENSOR_CLIFF_LEFT

    // The state of the cliff sensor on the front left of Roomba is sent as a 1
    // bit value (0 = no cliff, 1 = cliff).
    SENSOR_CLIFF_FRONT_LEFT

    // The state of the cliff sensor on the front right of Roomba is sent as a 1
    // bit value (0 = no cliff, 1 = cliff).
    SENSOR_CLIFF_FRONT_RIGHT

    // The state of the cliff sensor on the right side of Roomba is sent as a 1
    // bit value (0 = no cliff, 1 = cliff).
    SENSOR_CLIFF_RIGHT

    // The state of the virtual wall detector is sent as a 1 bit value
    // (0 = no virtual wall detected, 1 = virtual wall detected).
    SENSOR_VIRTUAL_WALL

    // The state of the four wheel overcurrent sensors are sent as individual
    // bits (0 = no overcurrent, 1 = overcurrent). There is no overcurrent
    // sensor for the vacuum on Roomba 500.
    SENSOR_WHEEL_OVERCURRENT

    _
    _

    // This value identifies the 8-bit IR character currently being received by
    // Roomba’s omnidirectional receiver.  A value of 0 indicates that no
    // character is being received. These characters include those sent by the
    // Roomba Remote, the Dock, Roomba 500 Virtual Walls, Create robots using
    // the Send-IR command, and user-created devices.
    SENSOR_IR_OMNI

    // The state of the Roomba buttons are sent as individual bits (0 = button
    // not pressed, 1 = button pressed). The day, hour, minute, clock, and
    // scheduling buttons that exist only on Roomba 560 and 570 will always
    // return 0 on a robot without these buttons.
    SENSOR_BUTTONS

    // The distance that Roomba has traveled in millimeters since the distance
    // it was last requested is sent as a signed 16-bit value, high byte first.
    // This is the same as the sum of the distance traveled by both wheels
    // divided by two. Positive values indicate travel in the forward direction;
    // negative values indicate travel in the reverse direction. If the value is
    // not polled frequently enough, it is capped at its minimum or maximum.
    // Range: -32768 – 32767
    SENSOR_DISTANCE

    // The angle in degrees that Roomba has turned since the angle was last
    // requested. Counter-clockwise angles are positive and clockwise angles
    // are negative. If the value is not polled frequently enough, it is capped
    // at its minimum or maximum. Range: -32768 – 32767
    SENSOR_ANGLE

    // This code indicates Roomba’s current charging state. Range: 0 – 5
    //
    //  Code Charging State
    //  0 Not charging
    //  1 Reconditioning Charging
    //  2 Full Charging
    //  3 Trickle Charging
    //  4 Waiting
    //  5 Charging Fault Condition
    SENSOR_CHARGING

    // This code indicates the voltage of Roomba’s battery in millivolts (mV).
    // Range: 0 – 65535 mV
    SENSOR_VOLTAGE

    // The current in milliamps (mA) flowing into or out of Roomba’s battery.
    // Negative currents indicate that the current is flowing out of the
    // battery, as during normal running. Positive currents indicate that the
    // current is flowing into the battery, as during charging.
    // Range: -32768 – 32767 mA
    SENSOR_CURRENT

    // The temperature of Roomba’s battery in degrees Celsius. Range: -128 – 127
    SENSOR_TEMPERATURE

    // The current charge of Roomba’s battery in milliamp-hours (mAh). The
    // charge value decreases as the battery is depleted during running and
    // increases when the battery is charged. Range: 0 – 65535 mAh
    SENSOR_BATTERY_CHARGE

    // The estimated charge capacity of Roomba’s battery in milliamp-hours (mAh). Range: 0 – 65535 mAh
    SENSOR_BATTERY_CAPACITY

    // The strength of the wall signal is returned as an unsigned 16-bit value.
    // Range: 0-1023.
    SENSOR_WALL_SIGNAL

    // The strength of the cliff left signal. Range: 0-4095.
    SENSOR_CLIFF_LEFT_SIGNAL

    // The strength of the cliff front left signal. Range: 0-4095.
    SENSOR_CLIFF_FRONT_LEFT_SIGNAL

    // The strength of the cliff front right signal. Range: 0-4095
    SENSOR_CLIFF_FRONT_RIGHT_SIGNAL

    // The strength of the cliff right signal. Range: 0-4095
    SENSOR_CLIFF_RIGHT_SIGNAL

    SENSOR_DIGITAL_INPUTS
    SENSOR_ANALOG_INPUT

    // Roomba’s connection to the Home Base and Internal Charger are returned as individual bits.
    SENSOR_CHARGING_SOURCE

    // The current OI mode is returned. Range 0-3.
    //  Number Mode
    //  0 Off
    //  1 Passive
    //  2 Safe
    //  3 Full
    SENSOR_OI_MODE

    // The currently selected OI song is returned. Range: 0-15
    SENSOR_SONG_NUMBER

    // The state of the OI song player is returned. 1 = OI song currently playing; 0 = OI song not playing.
    SENSOR_SONG_PLAYING

    // The number of data stream packets is returned.
    SENSOR_NUM_STREAM_PACKETS

    // The velocity most recently requested with a Drive command.
    SENSOR_REQUESTED_VELOCITY

    // The radius most recently requested with a Drive command.
    SENSOR_REQUESTED_RADIUS

    SENSOR_RIGHT_VELOCITY

    SENSOR_LEFT_VELOCITY
)

// SENSOR_PACKET_LENGTH is a map[SensorCode]byte that defines the length in bytes of sensor data packets.
var SENSOR_PACKET_LENGTH = map[SensorCode]byte{
    SENSOR_BUMP_WHEELS_DROPS:        1,
    SENSOR_WALL:                     1,
    SENSOR_CLIFF_LEFT:               1,
    SENSOR_CLIFF_FRONT_LEFT:         1,
    SENSOR_CLIFF_FRONT_RIGHT:        1,
    SENSOR_CLIFF_RIGHT:              1,
    SENSOR_VIRTUAL_WALL:             1,
    SENSOR_WHEEL_OVERCURRENT:        1,
    15:                              1,
    16:                              1,
    SENSOR_IR_OMNI:                  1,
    SENSOR_BUTTONS:                  1,
    SENSOR_DISTANCE:                 2,
    SENSOR_ANGLE:                    2,
    SENSOR_CHARGING:                 1,
    SENSOR_VOLTAGE:                  2,
    SENSOR_CURRENT:                  2,
    SENSOR_TEMPERATURE:              1,
    SENSOR_BATTERY_CHARGE:           2,
    SENSOR_BATTERY_CAPACITY:         2,
    SENSOR_WALL_SIGNAL:              2,
    SENSOR_CLIFF_LEFT_SIGNAL:        2,
    SENSOR_CLIFF_FRONT_LEFT_SIGNAL:  2,
    SENSOR_CLIFF_FRONT_RIGHT_SIGNAL: 2,
    SENSOR_CLIFF_RIGHT_SIGNAL:       2,
    SENSOR_DIGITAL_INPUTS:           1,
    SENSOR_ANALOG_INPUT:             2,
    SENSOR_CHARGING_SOURCE:          1,
    SENSOR_OI_MODE:                  1,
    SENSOR_SONG_NUMBER:              1,
    SENSOR_SONG_PLAYING:             1,
    SENSOR_NUM_STREAM_PACKETS:       1,
    SENSOR_REQUESTED_VELOCITY:       2,
    SENSOR_REQUESTED_RADIUS:         2,
    SENSOR_RIGHT_VELOCITY:           2,
    SENSOR_LEFT_VELOCITY:            2,
    0:                               26,
    1:                               10,
    2:                               6,
    3:                               10,
    4:                               14,
    5:                               12,
    6:                               52,
}
