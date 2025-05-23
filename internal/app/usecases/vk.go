package usecases

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/png"
	"log/slog"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/api/params"
	"github.com/SevereCloud/vksdk/v3/object"
	"github.com/ashurbekovz/vktexbot/internal/app/parsers"
	"github.com/ashurbekovz/vktexbot/internal/pkg/latex2img"
	"github.com/ashurbekovz/vktexbot/internal/pkg/template2img"
	"github.com/ashurbekovz/vktexbot/internal/tools/resize"
)

type VkOpt struct {
	PeerID       int
	Message      string
	Payload      string
	IsNewMessage bool
}

type VkRes struct {
}

type VkUsecase struct {
	vk     *api.VK
	t2i    *template2img.LatexTemplateToImgConverter
	parser *parsers.Vk
}

func NewVkUsecase(
	vk *api.VK,
	t2i *template2img.LatexTemplateToImgConverter,
	parser *parsers.Vk,
) *VkUsecase {
	return &VkUsecase{
		vk:     vk,
		t2i:    t2i,
		parser: parser,
	}
}

const (
	UnknownErrMessage = "Произошла неизвестная ошибка, мы уже выясняем причину. Возможно сработает повтор сообщения."
	EmptyErrMessage   = "Сгенерированное изображение оказалось пустым. Если это ошибка, пожалуйста, сообщите об этом разработчику."
)

func (u *VkUsecase) Execute(
	ctx context.Context,
	logger *slog.Logger,
	opt VkOpt,
) (VkRes, error) {
	messageParams, err := u.parser.Parse(opt.Message)
	if !messageParams.Mention && isGroup(opt.PeerID) {
		return VkRes{}, nil
	}

	logger.Info(
		"processing message",
		"message", opt.Message,
		"peer_id", opt.PeerID,
		"payload", opt.Payload,
		"is_new_message", opt.IsNewMessage,
		"is_group", isGroup(opt.PeerID),
	)

	if opt.Payload != "" && opt.Payload == `"help"` {
		b := params.NewMessagesSendBuilder()
		b.RandomID(0)
		b.PeerID(opt.PeerID)
		b.Message("Текст нужно вводить в LaTeX формате, заключая математические формулы в $ ... $, $$ ... $$ или их аналоги. Пример корректного сообщения:\n\nVkTeX $\\int_a^b f(x) \\, dx$ \n\nС подробным описанием функций можно ознакомиться по ссылке.")
		b.Keyboard(inlineHelpKeyboard())

		_, err := u.vk.MessagesSend(b.Params)
		if err != nil {
			u.sendErrToPeer(logger, opt.PeerID, UnknownErrMessage)
			return VkRes{}, fmt.Errorf("failed to send message: %w", err)
		}
		return VkRes{}, nil
	}

	if err != nil {
		u.sendErrToPeer(logger, opt.PeerID, UnknownErrMessage)
		return VkRes{}, fmt.Errorf("failed to parse message: %w", err)
	}

	if len(messageParams.Message) == 0 {
		u.sendErrToPeer(logger, opt.PeerID, EmptyErrMessage)
		return VkRes{}, nil
	}

	img, err := u.t2i.Convert(ctx, messageParams.Message, messageParams.ImageParams)
	if err != nil {
		var syntaxError *latex2img.SyntaxError
		var fullyTransparentError *resize.ImageFullyTransparentError
		switch {
		case errors.As(err, &syntaxError):
			u.sendErrToPeer(logger, opt.PeerID, syntaxError.UserError())
			return VkRes{}, nil
		case errors.As(err, &fullyTransparentError):
			u.sendErrToPeer(logger, opt.PeerID, EmptyErrMessage)
			return VkRes{}, nil
		default:
			u.sendErrToPeer(logger, opt.PeerID, UnknownErrMessage)
			return VkRes{}, fmt.Errorf("failed to convert message to image: %w", err)
		}
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		u.sendErrToPeer(logger, opt.PeerID, UnknownErrMessage)
		return VkRes{}, fmt.Errorf("failed to encode image: %v", err)
	}

	resp, err := u.vk.UploadMessagesPhoto(opt.PeerID, &buf)
	if err != nil {
		u.sendErrToPeer(logger, opt.PeerID, UnknownErrMessage)
		return VkRes{}, fmt.Errorf("failed to upload image: %w", err)
	}
	if len(resp) != 1 {
		u.sendErrToPeer(logger, opt.PeerID, UnknownErrMessage)
		return VkRes{}, fmt.Errorf("invalid response length from UploadMessagesPhoto: %d", len(resp))
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.PeerID(opt.PeerID)
	b.Attachment(resp[0].ToAttachment())
	if !isGroup(opt.PeerID) {
		b.Keyboard(helpKeyboard())
	}

	_, err = u.vk.MessagesSend(b.Params)
	if err != nil {
		u.sendErrToPeer(logger, opt.PeerID, UnknownErrMessage)
		return VkRes{}, fmt.Errorf("failed to send message: %w", err)
	}

	return VkRes{}, nil
}

func (u *VkUsecase) sendErrToPeer(
	logger *slog.Logger,
	peerID int,
	errMessage string,
) {
	logger.Info("attempting report to peer error", "err_message", errMessage)

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.PeerID(peerID)
	b.Message(errMessage)
	if !isGroup(peerID) {
		b.Keyboard(helpKeyboard())
	}
	_, err := u.vk.MessagesSend(b.Params)
	if err != nil {
		logger.Error("failed report to peer error", "error", err)
	}
}

func isGroup(peerID int) bool {
	return peerID >= 2000000000
}

func inlineHelpKeyboard() *object.MessagesKeyboard {
	keyboard := object.NewMessagesKeyboardInline()
	keyboard.AddRow()
	keyboard.AddOpenLinkButton("https://vk.com/@vktexbot-obnovlenie-bota-vktex-ot-04052020", "Список функций", "")
	return keyboard
}

func helpKeyboard() *object.MessagesKeyboard {
	keyboard := object.NewMessagesKeyboard(false)
	keyboard.AddRow()
	keyboard.AddTextButton("Помощь", "help", "primary")
	return keyboard
}
