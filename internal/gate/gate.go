// Package gate validates tick values against compiled game boundaries to prevent
// out-of-range temporal queries.
package gate

import "fmt"

// Validate checks that tick is a valid query index for the compiled game.
func Validate(tick int, maxTick uint16) error {
	if tick < 0 {
		return fmt.Errorf("tick %d is negative: must be >= 0", tick)
	}
	if tick > int(maxTick) {
		return fmt.Errorf("tick %d exceeds game length (%d elapsed seconds)", tick, maxTick)
	}
	return nil
}
