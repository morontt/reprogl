package models

const SubscriptionReplyComment = 1

type EmailSubscription struct {
	ID           int
	Email        string
	BlockSending bool
	Type         int
}
