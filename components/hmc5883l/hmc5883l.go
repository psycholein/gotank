package hmc5883l

import (
	"bytes"
	"encoding/binary"
	"gotank/i2c"
	"gotank/libs/embd"
	_ "gotank/libs/embd/host/all"
	"math"
)

const (
	HMC5883L_ADDR = 0x1E

	ConfigurationRegisterA  = 0x00
	ConfigurationRegisterB  = 0x01
	ModeRegister            = 0x02
	AxisDataRegister        = 0x03 // X,Z,Y (6 Byte)
	StatusRegister          = 0x09
	IdentificationRegisterA = 0x10
	IdentificationRegisterB = 0x11
	IdentificationRegisterC = 0x12
	MeasurementContinuous   = 0x00
	MeasurementSingleShot   = 0x01
	MeasurementIdle         = 0x03
)

type ThreeDData struct {
	X int16
	Z int16
	Y int16
}

type ThreeDDataFloat struct {
	X float64
	Y float64
	Z float64
}

type HMC5883LDriver struct {
	bus         embd.I2CBus
	compassRaw  ThreeDData
	Compass     ThreeDDataFloat
	initialized bool
	scale       float64
	declination float64
}

func NewHMC5883LDriver(bus embd.I2CBus) *HMC5883LDriver {
	m := &HMC5883LDriver{
		bus:         bus,
		initialized: false,
	}
	return m
}

func (h *HMC5883LDriver) read() (err error) {
	i2c.Lock()
	defer i2c.Unlock()
	if !h.initialized {
		h.initialized = true
		err = h.initialize()
		if err != nil {
			return
		}
	}
	ret := make([]byte, 6)
	err = h.bus.ReadFromReg(HMC5883L_ADDR, AxisDataRegister, ret)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(ret)

	binary.Read(buf, binary.BigEndian, &h.compassRaw)
	return nil
}

func (h *HMC5883LDriver) calculate() {
	if h.compassRaw.X == -4096 {
		h.Compass.X = 0
	} else {
		h.Compass.X = float64(h.compassRaw.X) * h.scale
	}

	if h.compassRaw.Y == -4096 {
		h.Compass.Y = 0
	} else {
		h.Compass.Y = float64(h.compassRaw.Y) * h.scale
	}

	if h.compassRaw.Z == -4096 {
		h.Compass.Z = 0
	} else {
		h.Compass.Z = -float64(h.compassRaw.Z) * h.scale
	}
}

func (h *HMC5883LDriver) Heading() float64 {
	h.read()
	h.calculate()
	headingRad := math.Atan2(h.Compass.Y, h.Compass.X) * 180 / math.Pi
	headingRad += h.declination + 180

	if headingRad < 0 {
		headingRad += 360
	}

	if headingRad > 360 {
		headingRad -= 360
	}

	return headingRad
}

func (h *HMC5883LDriver) write(reg byte, bytes []byte) (err error) {
	for _, b := range bytes {
		err = h.bus.WriteByteToReg(HMC5883L_ADDR, reg, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *HMC5883LDriver) initialize() (err error) {
	h.setScale(1.3)
	h.setContinuousMode()
	h.setDeclination(3, 28)

	return nil
}

func (h *HMC5883LDriver) setScale(gauss float64) {
	var scaleReg byte
	switch gauss {
	case 0.88:
		scaleReg = 0x00
		h.scale = 0.73
	case 1.3:
		scaleReg = 0x01
		h.scale = 0.92
	case 1.9:
		scaleReg = 0x02
		h.scale = 1.22
	case 2.5:
		scaleReg = 0x03
		h.scale = 1.52
	case 4.0:
		scaleReg = 0x04
		h.scale = 2.27
	case 4.7:
		scaleReg = 0x05
		h.scale = 3.03
	case 5.6:
		scaleReg = 0x06
		h.scale = 4.35
	case 8.1:
		scaleReg = 0x07
		h.scale = 0.73
	}
	scaleReg = scaleReg << 5
	h.setOptions(ConfigurationRegisterB, scaleReg)
}

func (h *HMC5883LDriver) setContinuousMode() {
	h.setOptions(ModeRegister, MeasurementContinuous)
}

func (h *HMC5883LDriver) setDeclination(degree int, min int) {
	h.declination = (float64(degree) + float64(min)/60)
}

func (h *HMC5883LDriver) setOptions(reg byte, o byte) (err error) {
	options := 0x00 | o
	err = h.write(reg, []byte{options})
	if err != nil {
		return
	}
	return nil
}
