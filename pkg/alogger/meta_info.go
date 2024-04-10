package alogger

import "fmt"

type MetaInfoKey int

const (
	LSIdKey MetaInfoKey = iota
	TPIdKey
	DocIdKey
	SenderIdKey
	ReceiverIdKey
)

// MetaInfo - Информация, выступающая в роли sub_trace_id. Нужна для обогащения событий дополнительной информацией о
// связанных ID ТП, ЛС, Документа, отправителя или получателя
type MetaInfo struct {
	LSId       string `json:"ls_id,omitempty"`
	TPId       string `json:"tp_id,omitempty"`
	DocId      string `json:"doc_id,omitempty"`
	SenderId   string `json:"sender_id,omitempty"`
	ReceiverId string `json:"receiver_id,omitempty"`
}

func (m *MetaInfo) IsEmpty() bool {
	return m.LSId == "" && m.TPId == "" && m.DocId == "" && m.SenderId == "" && m.ReceiverId == ""
}

func (m *MetaInfo) Bytes() []byte {
	if m.IsEmpty() {
		return nil
	}

	buffer := []byte("MetaInfo:\n")

	if m.LSId != "" {
		buffer = append(buffer, fmt.Sprintf("\t%s: %s\n", "ls_id", m.LSId)...)
	}

	if m.TPId != "" {
		buffer = append(buffer, fmt.Sprintf("\t%s: %s\n", "tp_id", m.TPId)...)
	}

	if m.DocId != "" {
		buffer = append(buffer, fmt.Sprintf("\t%s: %s\n", "doc_id", m.DocId)...)
	}

	if m.SenderId != "" {
		buffer = append(buffer, fmt.Sprintf("\t%s: %s\n", "sender_id", m.SenderId)...)
	}

	if m.ReceiverId != "" {
		buffer = append(buffer, fmt.Sprintf("\t%s: %s\n", "receiver_id", m.ReceiverId)...)
	}

	return buffer
}
