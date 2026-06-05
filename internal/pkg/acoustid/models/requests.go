package models

type AcoustIdLookupParams struct {
	Meta        []string `url:"meta" del:" "`
	Fingerprint string   `url:"fingerprint"`
	Duration    int      `url:"duration"`
}
