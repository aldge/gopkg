package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	epoch            = int64(1609459200000) // 自定义起始时间(2021-01-01 00:00:00 UTC)
	workerIDBits     = uint(5)              // 机器ID位数
	dataCenterIDBits = uint(5)              // 数据中心ID位数
	sequenceBits     = uint(12)             // 序列号位数

	maxWorkerID     = int64(-1) ^ (int64(-1) << workerIDBits)     // 最大机器ID
	maxDataCenterID = int64(-1) ^ (int64(-1) << dataCenterIDBits) // 最大数据中心ID
	maxSequence     = int64(-1) ^ (int64(-1) << sequenceBits)     // 最大序列号

	timeShift       = workerIDBits + dataCenterIDBits + sequenceBits // 时间戳左移位数
	dataCenterShift = sequenceBits + workerIDBits                    // 数据中心左移位数
	workerShift     = sequenceBits                                   // 机器ID左移位数
)

// Snowflake 结构体
type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	dataCenterID  int64
	sequence      int64
}

// NewSnowflake 创建一个Snowflake实例
func NewSnowflake(workerID, dataCenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, errors.New("worker ID out of range")
	}
	if dataCenterID < 0 || dataCenterID > maxDataCenterID {
		return nil, errors.New("data center ID out of range")
	}
	return &Snowflake{
		lastTimestamp: 0,
		workerID:      workerID,
		dataCenterID:  dataCenterID,
		sequence:      0,
	}, nil
}

// GenID 生成下一个唯一ID
func (s *Snowflake) GenID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1e6 // 当前毫秒时间戳

	if timestamp < s.lastTimestamp {
		return 0, errors.New("clock moved backwards")
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// 当前毫秒内序列号用完，等待下一毫秒
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	// 组合各部分生成ID
	id := ((timestamp - epoch) << timeShift) |
		(s.dataCenterID << dataCenterShift) |
		(s.workerID << workerShift) |
		s.sequence

	return id, nil
}
