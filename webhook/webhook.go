package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haakonleg/imt2681-assig2/mdb"
	"github.com/haakonleg/imt2681-assig2/router"
	"github.com/haakonleg/imt2681-assig2/ticker"
	"github.com/haakonleg/imt2681-assig2/util"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
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
	wh.db.Update(mdb.WEBHOOKS, nil, updateDoc)

	// Retrieve all webhooks where the counter is zero
	filter := bson.NewDocument(
		bson.EC.SubDocumentFromElements("triggerCount",
			bson.EC.Int64("$eq", 0)))

	webhooks := make([]*mdb.Webhook, 0)
	wh.db.Find(mdb.WEBHOOKS, filter, nil, &webhooks)

	// Invoke the webhooks
	for _, webhook := range webhooks {
		invokeWebhook(webhook, wh.db)

		// Reset the invoked webhook counter and set lastInvoked
		filter = bson.NewDocument(bson.EC.ObjectID("_id", webhook.ID))
		updateDoc = bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				bson.EC.Int64("triggerCount", webhook.MinTriggerValue),
				bson.EC.Int64("lastInvoked", util.NowMilli())))
		db.Update(mdb.WEBHOOKS, filter, updateDoc)
	}
}

// invokeWebhook sends a POST request to the webhook containing information about added tracks
func invokeWebhook(webhook *mdb.Webhook, db *mdb.Database) {
	// Build the request
	var request []byte
	ticker, er := ticker.MakeTicker(db, 0, webhook.LastInvoked)
	if er != nil {
		request, _ = json.Marshal(er)
	} else {
		request, _ = json.Marshal(ticker)
	}

	resp, err := http.Post(webhook.WebhookURL, "application/json", bytes.NewBuffer(request))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Invoked webhook %s. Status: %d", webhook.WebhookURL, resp.Status)
}

// GetWebhook is the handler for the API path GET /api/webhook/new_track/{webhook_id}
// Retrieves a webhook by the value of its ObjectID (hex encoded string)
func (wh *WebhookHandler) GetWebhook(req *router.Request) {
	webhookID := req.Vars["id"].(string)

	// Retrieve webhook from DB
	objectID, err := objectid.FromHex(webhookID)
	if err != nil {
		req.SendError(&router.Error{StatusCode: http.StatusBadRequest, Message: "Invalid ID"})
		return
	}
	filter := bson.NewDocument(bson.EC.ObjectID("_id", objectID))

	webhooks := make([]*mdb.Webhook, 0)
	if err := wh.db.Find(mdb.WEBHOOKS, filter, nil, &webhooks); err != nil {
		req.SendError(&router.Error{StatusCode: http.StatusInternalServerError, Message: "Internal database error"})
		return
	}
	if len(webhooks) < 1 {
		req.SendError(&router.Error{StatusCode: http.StatusBadRequest, Message: "Invalid ID"})
		return
	}

	req.SendJSON(&webhooks[0], http.StatusOK)
}

// DeleteWebhook is the handler for the API path DELETE /api/webhook/new_track/{webhook_id}
// Deletes a webhook by the value of its ObjectID (hex encoded string)
func (wh *WebhookHandler) DeleteWebhook(req *router.Request) {
	webhookID := req.Vars["id"].(string)

	// Delete webhook from DB
	objectID, err := objectid.FromHex(webhookID)
	if err != nil {
		req.SendError(&router.Error{StatusCode: http.StatusBadRequest, Message: "Invalid ID"})
		return
	}
	filter := bson.NewDocument(bson.EC.ObjectID("_id", objectID))

	delRes, err := wh.db.Delete(mdb.WEBHOOKS, filter)
	if err != nil {
		req.SendError(&router.Error{StatusCode: http.StatusInternalServerError, Message: "Internal database error"})
		return
	}
	if delRes.DeletedCount == 0 {
		req.SendError(&router.Error{StatusCode: http.StatusBadRequest, Message: "Invalid ID"})
		return
	}

	req.SendText("Webhook deleted", http.StatusOK)
}

// Register a webhook to be notified when new tracks are created
func (wh *WebhookHandler) PostWebhook(req *router.Request) {
	var webhookReq mdb.Webhook
	err := req.ParseJSONRequest(&webhookReq)
	if err != nil || len(webhookReq.WebhookURL) == 0 {
		req.SendError(&router.Error{StatusCode: http.StatusBadRequest, Message: "Invalid JSON"})
		return
	}

	webhook := mdb.CreateWebhook(webhookReq.WebhookURL, webhookReq.MinTriggerValue)
	id, err := wh.db.InsertObject(mdb.WEBHOOKS, &webhook)
	if err != nil {
		req.SendError(&router.Error{StatusCode: http.StatusInternalServerError, Message: "Internal database error"})
		return
	}

	response := struct {
		ID string `json:"id"`
	}{ID: id}
	req.SendJSON(&response, http.StatusOK)
}
