package job

import "log"

type dailyDecayProgressReporter struct {
	processed       uint32
	skipped, failed uint16
	errors          []error
}

func (rept *dailyDecayProgressReporter) AddProcessed() {
	rept.processed += 1
}

func (rept *dailyDecayProgressReporter) AddSkipped() {
	rept.skipped += 1
}

func (rept *dailyDecayProgressReporter) AddFailed(err error) {
	rept.failed += 1
	rept.errors = append(rept.errors, err)
}

func (rept *dailyDecayProgressReporter) Print() {
	log.Printf(
		"daily decay loop progress: processed=%d skipped=%d failed=%d",
		rept.processed,
		rept.skipped,
		rept.failed,
	)
}
