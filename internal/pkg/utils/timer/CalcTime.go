package timer

import "time"

type WrongInput struct {
	specification string
}

func (w *WrongInput) Error() string {
	return w.specification
}

func CalcTime(add int64, unit string) (int64, error) {
	now := time.Now().Unix()

	contrast := map[string]int64{
		"s": 1,
		"m": 60,
		"h": 3600,
		"d": 86400,
		"w": 604800,
		"M": 2419200,
		"y": 885427200,
	}

	if magnification, exist := contrast[unit]; exist != false {
		add = add * magnification
		return now + add, nil
	}

	return 0, &WrongInput{"Unknown unit"}
}
