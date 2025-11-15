package testsvr

import (
	"errors"
	"slices"
	"sync"
)

type DouglasPool struct {
	Firs  []*DouglasFir
	mutex sync.Mutex
}

func (pool *DouglasPool) OpenNewTest(uuid string, key string) error {
	fir, e := OpenOldTest(uuid, key)

	if e != nil {
		return e
	}

	pool.mutex.Lock()
	pool.Firs = append(pool.Firs, fir)
	pool.mutex.Unlock()

	go fir.OpenServer()

	return nil
}

func (pool *DouglasPool) CheckTestStatus(uuid string) bool {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	for _, v := range pool.Firs {
		if v.UUID == uuid {
			return true
		}
	}

	return false
}

func (pool *DouglasPool) CloseTest(uuid string) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	for i, v := range pool.Firs {
		if v.UUID == uuid {
			v.CloseServer()
			pool.Firs = slices.Delete(pool.Firs, i, i+1)
			return
		}
	}
}

func (pool *DouglasPool) GetServerIP(uuid string) string {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	for _, v := range pool.Firs {
		if v.UUID == uuid {
			return v.GetServerIP()
		}
	}

	return "NaN"
}

func (pool *DouglasPool) GetMark(testUUID string, candinateUUID string) (float64, error) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	for _, v := range pool.Firs {
		if v.UUID == testUUID {
			_, f, e := v.CalculateMark(candinateUUID)
			if e != nil {
				return 0, e
			}

			return f, nil
		}
	}

	return 0, errors.New("ERR_NO_SUCH_TEST")
}
