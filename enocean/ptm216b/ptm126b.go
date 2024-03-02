package ptm216b

import (
	"bytes"
	"encoding/binary"

	"github.com/go-ble/ble"
)

var enocean_identifier = []byte{0xda, 0x03} // reverse of what the doc says

type SwitchState byte

func (s SwitchState) IsButtonB1() bool {
	return s&(1<<4) > 0
}

func (s SwitchState) IsButtonB0() bool {
	return s&(1<<3) > 0
}

func (s SwitchState) IsButtonA1() bool {
	return s&(1<<2) > 0
}

func (s SwitchState) IsButtonA0() bool {
	return s&(1<<1) > 0
}

func (s SwitchState) IsPress() bool {
	return s&1 > 0
}

func (s SwitchState) IsRelease() bool {
	return s&1 == 0
}

type Event struct {
	Address  string
	State    SwitchState
	Sequence uint32
}

// NewEvent creates a ptm216b Event from a ble Advertisement
// If the Advertisement is not an enocean event, nil is returned and no error
func NewEvent(advertisement ble.Advertisement) (event *Event, err error) {
	rawManufacturerData := advertisement.ManufacturerData()
	if !bytes.Equal(rawManufacturerData[0:2], enocean_identifier) {
		return
	}
	event = &Event{
		Address:  advertisement.Addr().String(),
		Sequence: binary.LittleEndian.Uint32(rawManufacturerData[2:6]),
		State:    SwitchState(rawManufacturerData[6]),
	}
	return
}
