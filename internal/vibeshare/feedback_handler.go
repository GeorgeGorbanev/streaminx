package vibeshare

import (
	"fmt"
	"log/slog"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/templates"

	"github.com/tucnak/telebot"
)

func (vs *Vibeshare) FeedbackTextHandler(inMsg *telebot.Message) {
	slog.Info("handling feedback message",
		slog.String("from", inMsg.Sender.Username),
		slog.String("text", inMsg.Text))
	vs.feedbackRouter.RouteText(inMsg)
}

func (vs *Vibeshare) feedbackStart(inMsg *telebot.Message) {
	vs.respondToFeedback(&telegram.Message{To: inMsg.Sender, Text: templates.FeedbackStart})
}

func (vs *Vibeshare) feedback(inMsg *telebot.Message) {
	if err := vs.deliverFeedback(inMsg); err != nil {
		slog.Error("failed to deliver feedback", slog.Any("error", err))
	}

	vs.respondToFeedback(&telegram.Message{To: inMsg.Sender, Text: templates.FeedbackThanks})
}

func (vs *Vibeshare) feedbackReceiver() *telebot.User {
	return &telebot.User{ID: vs.feedbackReceiverID}
}

func (vs *Vibeshare) deliverFeedback(inMsg *telebot.Message) error {
	text := fmt.Sprintf(templates.FeedbackReceived, inMsg.Sender.Username, inMsg.Sender.ID, inMsg.Text)
	msg := telegram.Message{To: vs.feedbackReceiver(), Text: text}
	_, err := vs.feedbackSender.Send(&msg)
	return err
}

func (vs *Vibeshare) respondToFeedback(response *telegram.Message) {
	_, err := vs.feedbackSender.Send(response)
	if err != nil {
		slog.Error("failed to send message", slog.Any("error", err))
		return
	}
	slog.Info("sent feedback message response",
		slog.String("to", response.To.Recipient()),
		slog.String("text", response.Text))
}
