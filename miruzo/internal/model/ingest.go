package model

type IngestIDType = int64

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
