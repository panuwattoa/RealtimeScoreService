package utils

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Uint64ToString(number uint64) string {
	return strconv.FormatUint(number, 10)
}

func Int64ToString(number int64) string {
	return strconv.FormatInt(number, 10)
}

func ToInt64(s string) int64 {
	if s != "" {
		i64, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			zap.L().Error("cannot convert string to uint64", zap.Error(err))
		}
		return i64
	}
	return 0
}

func ToUint8(s string) uint8 {
	u64, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		zap.L().Error("cannot convert string to uint8", zap.Error(err))
	}
	return uint8(u64)
}

func ToUint16(s string) uint16 {
	u64, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		zap.L().Error("cannot convert string to uint16", zap.Error(err))
	}
	return uint16(u64)
}

func ToUint32(s string) uint32 {
	u64, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		zap.L().Error("cannot convert string to uint32", zap.Error(err))
	}
	return uint32(u64)
}

func ToUint64(s string) uint64 {
	if s != "" {
		u64, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			zap.L().Error("cannot convert string to uint64", zap.Error(err))
		}
		return u64
	}
	return 0
}

func ToFloat64(s string) float64 {
	u64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		zap.L().Error("cannot convert string to uint8", zap.Error(err))
	}
	return u64
}
func GetFormatDBTime() string {
	var timeFormat = "2006-01-02 15:04:05"
	return timeFormat
}

func TimeToString(dateTime time.Time) string {
	timeString := dateTime.Format("2006-01-02 15:04:05")
	return timeString
}

func GetTimeNowString(local *time.Location) string {
	timeString := time.Now().In(local).Format(GetFormatDBTime())
	return timeString
}

func GetCurrentUnixTimestampString() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func GetTimeFromUnixString(unixString string) time.Time {
	return time.Unix(ToInt64(unixString), 0)
}
