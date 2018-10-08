package paragliding

import "github.com/mongodb/mongo-go-driver/bson/objectid"

// Webhook is the model of registered webhooks in the database
// MinTriggerValue is the limit to how many new tracks are created before the webhook is notified
// TriggerCount is decremented by one each time a new track is created, to know when to notify (when it is 0)
// LastInvoked is a timestamp of when the webhook was last invoked
type Webhook struct {
	ID              objectid.ObjectID `bson:"_id" json:"-"`
	WebhookURL      string            `bson:"webhookURL" json:"webhookURL"`
	MinTriggerValue int64             `bson:"minTriggerValue" json:"minTriggerValue"`
	TriggerCount    int64             `bson:"triggerCount" json:"-"`
	LastInvoked     int64             `bson:"lastInvoked" "json:"-"`
}

func createWebhook(webhookUrl string, minTriggerValue int64) Webhook {
	// If minTriggerValue was not specified, set to 1
	if minTriggerValue == 0 {
		minTriggerValue = 1
	}

	return Webhook{
		ID:              objectid.New(),
		WebhookURL:      webhookUrl,
		MinTriggerValue: minTriggerValue,
		TriggerCount:    minTriggerValue,
		LastInvoked:     -1}
}
