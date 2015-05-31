package l293d

import (
	"gotank/libs/embd"
	_ "gotank/libs/embd/host/all"
)

type MotorShieldL293d struct {
	latch, clk, enable, data, pwm, motor int
}

const (
	motor1A = 2
	motor1B = 3
	motor2A = 1
	motor2B = 4
	motor4A = 0
	motor4B = 6
	motor3A = 5
	motor3B = 7

	forward  = 1
	backward = 2
	release  = 3
)

var latchState byte = 0

func InitMotor(latch int, clk int, enable int, data int, pwm int, motor int) *MotorShieldL293d {
	m := MotorShieldL293d{latch, clk, enable, data, pwm, motor}
	m.init()
	return &m
}

func (m MotorShieldL293d) init() {
	embd.SetDirection(m.latch, embd.Out)
	embd.SetDirection(m.enable, embd.Out)
	embd.SetDirection(m.data, embd.Out)
	embd.SetDirection(m.clk, embd.Out)
	m.latchTx()
	embd.DigitalWrite(m.enable, embd.Low)
}

func (m MotorShieldL293d) initMotor() {
	if m.motor == 1 {
		latchState &= m.bv(motor1A) & m.bv(motor1B)
	}
	if m.motor == 2 {
		latchState &= m.bv(motor2A) & m.bv(motor2B)
	}
	if m.motor == 3 {
		latchState &= m.bv(motor3A) & m.bv(motor3B)
	}
	if m.motor == 4 {
		latchState &= m.bv(motor4A) & m.bv(motor4B)
	}
	m.latchTx()
	m.Speed(1023)
}

func (m MotorShieldL293d) command(cmd int) {
	var a, b byte
	if m.motor == 1 {
		a = motor1A
		b = motor1B
	}
	if m.motor == 2 {
		a = motor2A
		b = motor2B
	}
	if m.motor == 3 {
		a = motor3A
		b = motor3B
	}
	if m.motor == 4 {
		a = motor4A
		b = motor4B
	}
	if cmd == forward {
		latchState |= m.bv(a)
		latchState &= m.bv(b) ^ 255
	}
	if cmd == backward {
		latchState &= m.bv(a) ^ 255
		latchState |= m.bv(b)
	}
	if cmd == release {
		latchState &= m.bv(a) ^ 255
		latchState &= m.bv(b) ^ 255
	}
	m.latchTx()
}

func (m MotorShieldL293d) latchTx() {
	embd.DigitalWrite(m.latch, embd.Low)
	embd.DigitalWrite(m.data, embd.Low)
	var i byte
	for i = 0; i < 8; i++ {
		embd.DigitalWrite(m.clk, embd.Low)
		if latchState&m.bv(7-i) > 0 {
			embd.DigitalWrite(m.data, embd.High)
		} else {
			embd.DigitalWrite(m.data, embd.Low)
		}
		embd.DigitalWrite(m.clk, embd.High)
	}
	embd.DigitalWrite(m.latch, embd.High)
}

func (m MotorShieldL293d) bv(i byte) byte {
	return 1 << i
}

func (m MotorShieldL293d) Forward() {
}

func (m MotorShieldL293d) Backward() {
}

func (m MotorShieldL293d) Speed(i int) {

}

func (m MotorShieldL293d) Stop() {
}
