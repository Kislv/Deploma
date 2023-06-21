package cast

import (
	"encoding/binary"
	"fmt"
	"read-adviser-bot/utils/log"
	"math"
	"runtime"
	"strconv"
	"time"

	"github.com/jackc/pgx/pgtype"
)

func ToString(src []byte) string {
	return string(src)
}

func IntToStr(src uint64) string {
	return fmt.Sprint(src)
}

func Uint16ToStr(src uint16) string {
	return fmt.Sprint(src)
}

func FlToStr(src float64) string {
	return fmt.Sprintf("%.1f", src)
}

func TimeToStr(src time.Time, withTime bool) string {
	if withTime {
		return src.Format("2006.01.02 15:04:05")
	}
	return src.Format("2006.01.02")
}

func StringToDate(src string) (time.Time, error){
	var result time.Time
	year, err := strconv.Atoi(src[0:4]) 
	if err != nil {
		return result, err
	}
	month, err := strconv.Atoi(src[5:7])
	if err != nil {
		return result, err
	}
	day, err := strconv.Atoi(src[8:])
	if err != nil {
		return result, err
	}
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}

func ToUint64(src []byte) uint64 {
	return binary.BigEndian.Uint64(src)
}

func ToUint32(src []byte) uint32 {
	return binary.BigEndian.Uint32(src)
}

func ToUint16(src []byte) uint16 {
	return binary.BigEndian.Uint16(src)
}

func ToUint8(src []byte) uint8 {
	return uint8(binary.BigEndian.Uint16(src))
}

func ToFloat64(src []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(src))
}

func StringToFloat64(src string) (float64, error) {
	s, err := strconv.ParseFloat(src, 64)
	if err != nil {
		log.Error(err)
		return 0.0, err
	}
	return s, nil
}

func ToTime(src []byte) time.Time {
	tmp := pgtype.Timestamp{}
	err := tmp.DecodeBinary(nil, src)
	if err != nil {
		log.Error(err)
		return time.Time{}
	}
	return tmp.Time
}

func ToDate(src []byte) time.Time {
	tmp := pgtype.Date{}
	err := tmp.DecodeBinary(nil, src)
	if err != nil {
		log.Error(err)
		return time.Time{}
	}
	return tmp.Time
}

func ToBool(src []byte) bool {
	tmp := pgtype.Bool{}
	err := tmp.DecodeBinary(nil, src)
	if err != nil {
		log.Error(err)
		return tmp.Bool
	}
	return tmp.Bool
}
func GetCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
