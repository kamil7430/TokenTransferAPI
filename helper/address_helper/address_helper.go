package address_helper

import (
	"encoding/hex"
	"errors"
	"strings"
)

const addressLength = 42

func CheckAddress(address string) error {
	if len(address) != addressLength {
		return errors.New("invalid address length")
	}

	if !strings.HasPrefix(address, "0x") {
		return errors.New("address must start with 0x")
	}

	addressWithout0x := address[2:]
	dst := make([]byte, hex.DecodedLen(len(addressWithout0x)))
	_, err := hex.Decode(dst, []byte(addressWithout0x))
	if err != nil {
		return errors.New("invalid hex number")
	}

	return nil
}
