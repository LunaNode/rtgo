package rtgo

import "time"

type Ticket struct {
	ID string `rt:"id"`
	Queue string
	Owner string
	Creator string
	Subject string
	Status string
	Priority int
	InitialPriority int
	FinalPriority int
	Requestors string
	Cc string
	AdminCc string
	Created time.Time
	Starts time.Time
	Started time.Time
	Due time.Time
	Resolved time.Time
	Told time.Time
	LastUpdated time.Time
	TimeEstimated int
	TimeWorked int
	TimeLeft int
}
