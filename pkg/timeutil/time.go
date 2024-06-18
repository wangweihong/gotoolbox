package timeutil

import (
	"fmt"
	"strings"
	"time"
)

type SdkTime time.Time

func (t *SdkTime) UnmarshalJSON(data []byte) error {
	tmp := strings.Trim(string(data[:]), "\"")

	now, err := time.ParseInLocation(`2006-01-02T15:04:05Z`, tmp, time.UTC)
	if err == nil {
		*t = SdkTime(now)
		return err
	}

	now, err = time.ParseInLocation(`2006-01-02T15:04:05`, tmp, time.UTC)
	if err == nil {
		*t = SdkTime(now)
		return err
	}

	now, err = time.ParseInLocation(`2006-01-02 15:04:05`, tmp, time.UTC)
	if err == nil {
		*t = SdkTime(now)
		return err
	}

	now, err = time.ParseInLocation(`2006-01-02T15:04:05+08:00`, tmp, time.UTC)
	if err == nil {
		*t = SdkTime(now)
		return err
	}

	now, err = time.ParseInLocation(time.RFC3339, tmp, time.UTC)
	if err == nil {
		*t = SdkTime(now)
		return err
	}

	now, err = time.ParseInLocation(time.RFC3339Nano, tmp, time.UTC)
	if err == nil {
		*t = SdkTime(now)
		return err
	}

	return err
}

func (t SdkTime) MarshalJSON() ([]byte, error) {
	rs := []byte(fmt.Sprintf(`"%s"`, t.String()))
	return rs, nil
}

func (t SdkTime) String() string {
	return time.Time(t).Format(`2006-01-02T15:04:05Z`)
}

func ParseTime(data string) (time.Time, error) {
	tmp := strings.Trim(data[:], "\"")

	now, err := time.ParseInLocation(`2006-01-02T15:04:05Z`, tmp, time.UTC)
	if err == nil {
		return now, nil
	}

	now, err = time.ParseInLocation(`2006-01-02T15:04:05`, tmp, time.UTC)
	if err == nil {
		return now, nil
	}

	now, err = time.ParseInLocation(`2006-01-02 15:04:05`, tmp, time.UTC)
	if err == nil {
		return now, nil
	}

	now, err = time.ParseInLocation(`2006-01-02T15:04:05+08:00`, tmp, time.UTC)
	if err == nil {
		return now, nil
	}

	now, err = time.ParseInLocation(time.RFC3339, tmp, time.UTC)
	if err == nil {
		return now, nil
	}

	now, err = time.ParseInLocation(time.RFC3339Nano, tmp, time.UTC)
	if err == nil {
		return now, nil
	}

	return time.Time{}, err
}

func FormatDuration(duration time.Duration) string {
	seconds := int(duration.Seconds())
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	remainingSeconds := seconds % 60

	var result string

	if hours > 0 {
		result += fmt.Sprintf("%d小时", hours)
	}
	if minutes > 0 || hours > 0 {
		result += fmt.Sprintf("%d分", minutes)
	}
	result += fmt.Sprintf("%d秒", remainingSeconds)

	return result
}
