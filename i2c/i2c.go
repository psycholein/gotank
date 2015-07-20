package i2c

import "sync"

var i2cMutex = &sync.Mutex{}

func Lock() {
	i2cMutex.Lock()
}

func Unlock() {
	i2cMutex.Unlock()

}
