package roller

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func ParseRollString(s string) (int64, string, error) {
	s = strings.ReplaceAll(s, "-", "+-")
	parts := strings.Split(s, "+")
	var out int64
	var outString string = " <= "
	for i, part := range parts {
		adv := false
		dis := false
		if strings.HasSuffix(part, "a") {
			adv = true
		} else if strings.HasSuffix(part, "d") {
			dis = true
		}
		part = strings.TrimSuffix(part, "a")
		part = strings.TrimSuffix(part, "d")
		if c, err := strconv.ParseInt(part, 10, 64); err == nil && i != 0 {
			out += c
			if c < 0 {
				outString += " - " + strconv.FormatInt(-1*c, 10)
			} else {
				outString += " + " + strconv.FormatInt(c, 10)
			}
		} else {
			i, s, err := processDie(part, adv, dis)
			if err != nil {
				return i, outString, err
			}
			out += i
			outString += s
		}
	}
	return out, strings.Replace(outString, " <=  + ", " <= ", 1), nil
}

func processDie(part string, adv, dis bool) (int64, string, error) {
	var mult int64 = 1
	var size int64
	var outString string = ""
	var val int64
	if strings.HasPrefix(part, "-") {
		mult = -1
		part = strings.TrimPrefix(part, "-")
	}
	if strings.Contains(part, "d") {
		ps := strings.Split(part, "d")
		count, err := strconv.Atoi(ps[0])
		if err != nil {
			count = 1
		}
		size, err = strconv.ParseInt(ps[1], 10, 64)
		if err != nil {
			return -1, outString, errors.New("invalid format")
		}
		for i := 0; i < count; i++ {
			i, s, err := rollDie(size, adv, dis)
			if err != nil {
				return i, outString, err
			}
			val += i * mult
			if mult > 0 {
				outString += " + " + s
			} else if mult < 0 {
				outString += " - " + s
			} else if mult > 0 {
				outString += s
			}
		}
	} else if strings.Contains(part, "x") {
		ps := strings.Split(part, "x")
		if len(ps) < 2 {
			return -1, outString, errors.New("invalid format")
		}
		modifer, err := strconv.ParseInt(ps[0], 10, 64)
		if err != nil {
			return -1, outString, errors.New("invalid format")
		}
		multiple, err := strconv.ParseInt(ps[1], 10, 64)
		if err != nil {
			return -1, outString, errors.New("invalid format")
		}
		out := modifer * multiple * mult
		if out < 0 {
			outString += " - " + strconv.FormatInt(-1*out, 10)
		} else {
			outString += " + " + strconv.FormatInt(out, 10)
		}
		val += out
	} else {
		size, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return -1, outString, err
		}
		i, s, err := rollDie(size, adv, dis)
		if err != nil {
			return i, outString, err
		}
		val += i * mult
		if mult > 0 {
			outString += " + " + s
		} else if mult < 0 {
			outString += " - " + s
		} else if mult > 0 {
			outString += s
		}
	}
	return val, outString, nil
}

func rollDie(i int64, adv, dis bool) (int64, string, error) {
	if i <= 1 {
		return -2, "", errors.New("invalid die")
	}
	size := big.NewInt(i)
	one := big.NewInt(1)
	num, err := rand.Int(rand.Reader, size)
	if err != nil {
		return -1, "", errors.New("can't get die roll")
	}
	num.Add(num, one)
	if i == 10 && num.Cmp(big.NewInt(10)) == 0 {
		num = big.NewInt(0)
	}
	if adv {
		num2, err := rand.Int(rand.Reader, size)
		if err != nil {
			return -1, "", errors.New("can't get die roll")
		}
		num2.Add(num2, one)
		if num2.Cmp(num) == 1 {
			return num2.Int64(), fmt.Sprintf("{ ~~%s~~ %s }", num.String(), num2.String()), nil
		}
		return num.Int64(), fmt.Sprintf("{ ~~%s~~ %s }", num2.String(), num.String()), nil
	} else if dis {
		num2, err := rand.Int(rand.Reader, size)
		if err != nil {
			return -1, "", errors.New("can't get die roll")
		}
		num2.Add(num2, one)
		if num2.Cmp(num) == -1 {
			return num2.Int64(), fmt.Sprintf("{ ~~%s~~ %s }", num.String(), num2.String()), nil
		}
		return num.Int64(), fmt.Sprintf("{ ~~%s~~ %s }", num2.String(), num.String()), nil
	}
	return num.Int64(), fmt.Sprintf("{ %s }", num.String()), nil
}
