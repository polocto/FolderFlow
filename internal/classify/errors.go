package classify

import "fmt"

func ErrInvalidRegroupMode(mode string) error {
	return fmt.Errorf("invalid regroup mode: %s", mode)
}
