package auth

import "strconv"

func CheckDonation(amount string) (int, error) {
	number, err := strconv.Atoi(string(amount))
	if err != nil {
		return 0, err
	}

	return number, nil
}
