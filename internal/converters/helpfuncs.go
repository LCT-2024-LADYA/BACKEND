package converters

import (
	"gopkg.in/guregu/null.v3"
	"strconv"
	"time"
)

func getStringPointer(s null.String) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func getNullString(s *string) null.String {
	if s == nil {
		return null.NewString("", false)
	}
	return null.NewString(*s, true)
}

func getIntPointer(i null.Int) *int {
	if i.Valid {
		a := int(i.Int64)
		return &a
	}
	return nil
}

func getNullInt(i *int) null.Int {
	if i == nil {
		return null.NewInt(0, false)
	}
	return null.NewInt(int64(*i), true)
}

func getStringPointerCursor(i int) *string {
	if i == -1 {
		return nil
	}
	a := strconv.Itoa(i)
	return &a
}

func getIntCursor(s string) (int, error) {
	if s == "" {
		return 0, nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func getNullBool(b *bool) null.Bool {
	if b == nil {
		return null.NewBool(false, false)
	}
	return null.NewBool(*b, true)
}

func getNullTime(t *time.Time) null.Time {
	if t == nil {
		return null.NewTime(time.Time{}, false)
	}
	return null.NewTime(*t, true)
}
