package checks

import (
	"fmt"
	"net/url"
	"time"
)

func CheckUrl(s string) error {
	if s == "" {
		return fmt.Errorf("url not specified")
	}

	parsed, err := url.Parse(s)
	if err != nil {
		return err
	}

	if parsed.Scheme == "" {
		return fmt.Errorf("scheme not specified")
	}
	return nil
}

func CheckPeriod(p time.Duration) error {
	if p <= 0 {
		return fmt.Errorf("period can't be less or equal to zero %v", p)
	}
	return nil
}
