package models

type EmailSubscription struct {
	ID           int
	Email        string
	BlockSending bool
}
