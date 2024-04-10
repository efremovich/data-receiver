package broker

import "time"

type Task struct {
	Description []byte    `json:"description"`
	IsReceipt   bool      `json:"isReceipt"`
	ReceiptURL  string    `json:"receiptURL"`
	Created     time.Time `json:"created"`
}

const ReceiverStreamName = "receiver_outgoing"

const ReceiverSubjectNormalPriority = "receiver_outgoing.normal_priority"
