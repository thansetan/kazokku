package helpers

func GetLast4Digits(ccNumber string) string {
	if len(ccNumber) < 4 {
		return ccNumber
	}
	return ccNumber[len(ccNumber)-4:]
}
