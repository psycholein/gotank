package mpu6050

import (
	"bytes"
	"encoding/binary"
	"gotank/i2c"
	"gotank/libs/embd"
	_ "gotank/libs/embd/host/all"
)

const MPU6050_RA_ACCEL_XOUT_H = 0x3B
const MPU6050_RA_PWR_MGMT_1 = 0x6B
const MPU6050_PWR1_CLKSEL_BIT = 2
const MPU6050_PWR1_CLKSEL_LENGTH = 3
const MPU6050_CLOCK_PLL_XGYRO = 0x01
const MPU6050_GYRO_FS_250 = 0x00
const MPU6050_RA_GYRO_CONFIG = 0x1B
const MPU6050_GCONFIG_FS_SEL_LENGTH = 2
const MPU6050_GCONFIG_FS_SEL_BIT = 4
const MPU6050_RA_ACCEL_CONFIG = 0x1C
const MPU6050_ACONFIG_AFS_SEL_BIT = 4
const MPU6050_ACONFIG_AFS_SEL_LENGTH = 2
const MPU6050_ACCEL_FS_2 = 0x00
const MPU6050_PWR1_SLEEP_BIT = 6
const MPU6050_ADDR = 0x68

type ThreeDData struct {
	X int16
	Y int16
	Z int16
}

type MPU6050Driver struct {
	bus           embd.I2CBus
	Accelerometer ThreeDData
	Gyroscope     ThreeDData
	Temperature   float64
	initialized   bool
}

func NewMPU6050Driver(bus embd.I2CBus) *MPU6050Driver {
	m := &MPU6050Driver{
		bus:         bus,
		initialized: false,
	}
	return m
}

func (h *MPU6050Driver) Read() (err error) {
	i2c.Lock()
	defer i2c.Unlock()
	if !h.initialized {
		h.initialized = true
		err = h.initialize()
		if err != nil {
			return
		}
	}
	ret := make([]byte, 14)
	err = h.bus.ReadFromReg(MPU6050_ADDR, MPU6050_RA_ACCEL_XOUT_H, ret)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(ret)

	var temperature int16

	binary.Read(buf, binary.BigEndian, &h.Accelerometer)
	binary.Read(buf, binary.BigEndian, &temperature)
	binary.Read(buf, binary.BigEndian, &h.Gyroscope)
	h.Temperature = float64(temperature)/340 + 36.53
	return nil
}

func (h *MPU6050Driver) GetGyroscope() (x int16, y int16, z int16) {
	return h.Gyroscope.X, h.Gyroscope.Y, h.Gyroscope.Z
}

func (h *MPU6050Driver) GetAccelerometer() (x int16, y int16, z int16) {
	return h.Accelerometer.X, h.Accelerometer.Y, h.Accelerometer.Z
}

func (h *MPU6050Driver) GetTemperature() (celsius float64) {
	return h.Temperature
}

func (h *MPU6050Driver) write(reg byte, bytes []byte) (err error) {
	for _, b := range bytes {
		err = h.bus.WriteByteToReg(MPU6050_ADDR, reg, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *MPU6050Driver) initialize() (err error) {
	err = h.write(MPU6050_RA_PWR_MGMT_1, []byte{
		MPU6050_PWR1_CLKSEL_BIT,
		MPU6050_PWR1_CLKSEL_LENGTH,
		MPU6050_CLOCK_PLL_XGYRO})
	if err != nil {
		return
	}

	err = h.write(MPU6050_GYRO_FS_250, []byte{
		MPU6050_RA_GYRO_CONFIG,
		MPU6050_GCONFIG_FS_SEL_LENGTH,
		MPU6050_GCONFIG_FS_SEL_BIT})
	if err != nil {
		return
	}

	err = h.write(MPU6050_RA_ACCEL_CONFIG, []byte{
		MPU6050_ACONFIG_AFS_SEL_BIT,
		MPU6050_ACONFIG_AFS_SEL_LENGTH,
		MPU6050_ACCEL_FS_2})
	if err != nil {
		return
	}

	err = h.write(MPU6050_RA_PWR_MGMT_1, []byte{
		MPU6050_PWR1_SLEEP_BIT, 0})
	if err != nil {
		return
	}

	return nil
}
