package mode

import (
	"errors"
	"regexp"
)

func checkPort(port string) error {
	re, err := regexp.Compile(`\D`)
	if err != nil {
		return err
	}

	if len(re.FindString(port)) != 0 {
		return errors.New("port params error")
	}
	return nil
}

func checkAddress(address string) error {
	re, err := regexp.Compile(`[^0-9a-zA-Z:.]`)
	if err != nil {
		return err
	}

	if re.MatchString(address) {
		return errors.New("address params error")
	}
	return nil
}
