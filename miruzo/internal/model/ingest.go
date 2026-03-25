package model

type IngestIDType = int64

const (
	// MinIngestID is the minimum valid ingest identifier.
	MinIngestID IngestIDType = 1
	// MaxIngestID is the maximum valid ingest identifier (2^53 - 1).
	// This keeps IDs within JavaScript's safe integer range.
	MaxIngestID IngestIDType = 9007199254740991
)

type ProcessStatus uint8

const (
	ProcessStatusProcessing ProcessStatus = iota
	ProcessStatusFinished
)

type VisibilityStatus uint8

const (
	VisibilityStatusPrivate = iota
	VisibilityStatusPublic
)
