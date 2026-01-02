// Copyright 2026 Paul Sade
// GPLv3 - See LICENSE for details.


package classify

import "fmt"

func ErrInvalidRegroupMode(mode string) error {
	return fmt.Errorf("invalid regroup mode: %s", mode)
}
