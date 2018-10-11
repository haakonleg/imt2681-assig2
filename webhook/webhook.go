package webhook

import (
	"net/http"

	"github.com/haakonleg/imt2681-assig2/router"
	"github.com/mongodb/mongo-go-driver/bson"
)

type WebhookHandler struct {
	db *mdb.Database
}

func NewWebhookHandler(db *mdb.Database) *WebhookHandler {
	return &WebhookHandler{
		db: db}
}

// CheckInvokeWebhooks decrements the trigger counters of each webhook by one, then checks which webhooks that have their counter/trigger
// equal to zero and invokes the ones who have, then their counter is reset
func (wh *WebhookHandler) CheckInvokeWebhooks(db *mdb.Database) {
	// Decrement all webhooks triggercount by one
	updateDoc := bson.NewDocument(
		bson.EC.SubDocumentFromElements("$inc",
			bson.EC.Int64("triggerCount", -1)))
	db.Update(WEBHOOKS, nil, updateDoc)

	// Retrieve all webhooks where the counter is zero
	filter := bson.NewDocument(
		bson.EC.SubDocumentFromElements("triggerCount",
			bson.EC.Int64("$eq", 0)))
	webhooks, _ := db.Find(WEBHOOKS, filter, nil)

	// Invoke the webhooks
	for _, webhook := range webhooks {
		invokeWebhook(webhook.(Webhook), db)

		// Reset the invoked webhook counter and set lastInvoked
		filter = bson.NewDocument(bson.EC.ObjectID("_id", webhook.(Webhook).ID))
		updateDoc = bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				bson.EC.Int64("triggerCount", webhook.(Webhook).MinTriggerValue),
				bson.EC.Int64("lastInvoked", nowMilli())))
		db.Update(WEBHOOKS, filter, updateDoc)
	}
}

// invokeWebhook sends a POST request to the webhook containing information about added tracks
func invokeWebhook(webhook Webhook, db *mdb.Database) {
	//req, err := http.NewRequest(http.MethodPost, webhook.WebhookURL)
}

// TODO: Implement
func getWebhook(req *router.Request) {

}

// TODO: Implement
func deleteWebhook(req *router.Request) {

}

// Register a webhook to be notified when new tracks are created
func (wh *WebhookHandler) PostWebhook(req *router.Request) {
	var webhookReq Webhook
	err := req.ParseJSONRequest(&webhookReq)
	if err != nil || len(webhookReq.WebhookURL) == 0 {
		req.SendError("Invalid request", http.StatusBadRequest)
		return
	}

	webhook := createWebhook(webhookReq.WebhookURL, webhookReq.MinTriggerValue)
	id, err := wh.db.InsertObject(WEBHOOKS, &webhook)
	if err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}

	response := struct {
		ID string `json:"id"`
	}{ID: id}
	req.SendJSON(&response, http.StatusOK)
}
