package statuses

type Status int

const (
	Created Status = iota
	Enqueued
	Calculating
	Finished
	Failed
)
