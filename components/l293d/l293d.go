package l293d

import (
	"gotank/libs/embd"
	_ "gotank/libs/embd/host/all"
	"time"
)

type MotorL293dInterface interface {
	Forward()
	Backward()
	Speed(i int)
	Stop()
	Release()
}

type MotorL293d struct {
	pwm, motor int
	l293d      *L293d
	pin        embd.PWMPin
}

type L293d struct {
	latch, clk, enable, data int
	latchState               byte
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

func bv(i byte) byte {
	return 1 << i
}

func InitL293d(latch int, clk int, enable int, data int) L293d {
	l293d := L293d{latch, clk, enable, data, 0}
	l293d.init()
	return l293d
}

func (l *L293d) init() {
	embd.SetDirection(l.latch, embd.Out)
	embd.SetDirection(l.enable, embd.Out)
	embd.SetDirection(l.data, embd.Out)
	embd.SetDirection(l.clk, embd.Out)
	l.latchTx()
	embd.DigitalWrite(l.enable, embd.Low)
}

func (l *L293d) InitMotor(pwm int, motor int) *MotorL293d {
	m := &MotorL293d{pwm, motor, l, nil}
	m.init()
	return m
}

func (l *L293d) latchTx() {
	embd.DigitalWrite(l.latch, embd.Low)
	time.Sleep(1 * time.Microsecond)
	embd.DigitalWrite(l.data, embd.Low)
	time.Sleep(1 * time.Microsecond)
	var i byte
	for i = 0; i < 8; i++ {
		embd.DigitalWrite(l.clk, embd.Low)
		time.Sleep(1 * time.Microsecond)
		if l.latchState&bv(7-i) > 0 {
			embd.DigitalWrite(l.data, embd.High)
		} else {
			embd.DigitalWrite(l.data, embd.Low)
		}
		time.Sleep(1 * time.Microsecond)
		embd.DigitalWrite(l.clk, embd.High)
		time.Sleep(1 * time.Microsecond)
	}
	time.Sleep(1 * time.Microsecond)
	embd.DigitalWrite(l.latch, embd.High)
}

func (m MotorL293d) init() {
	m.pin, _ = embd.NewPWMPin(m.pwm)
}

func (m MotorL293d) initMotor() {
	if m.motor == 1 {
		m.l293d.latchState &= bv(motor1A) ^ 255&bv(motor1B) ^ 255
	}
	if m.motor == 2 {
		m.l293d.latchState &= bv(motor2A) ^ 255&bv(motor2B) ^ 255
	}
	if m.motor == 3 {
		m.l293d.latchState &= bv(motor3A) ^ 255&bv(motor3B) ^ 255
	}
	if m.motor == 4 {
		m.l293d.latchState &= bv(motor4A) ^ 255&bv(motor4B) ^ 255
	}
	m.l293d.latchTx()
	m.Speed(100)
}

func (m MotorL293d) command(cmd int) {
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
		m.l293d.latchState |= bv(a)
		m.l293d.latchState &= bv(b) ^ 255
	}
	if cmd == backward {
		m.l293d.latchState &= bv(a) ^ 255
		m.l293d.latchState |= bv(b)
	}
	if cmd == release {
		m.l293d.latchState &= bv(a) ^ 255
		m.l293d.latchState &= bv(b) ^ 255
	}
	m.l293d.latchTx()
}

func (m MotorL293d) Forward() {
	m.command(forward)
}

func (m MotorL293d) Backward() {
	m.command(backward)
}

func (m MotorL293d) Speed(i int) {
	if m.pin != nil {
		m.pin.SetDuty(i)
	}
}

func (m MotorL293d) Stop() {
	m.command(release)
}

func (m MotorL293d) Release() {
	m.pin.Close()
}
