package usecases

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"log"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/api/params"
	"github.com/ashurbekovz/vktexbot/internal/app/parsers"
	"github.com/ashurbekovz/vktexbot/internal/pkg/template2img"
)

type VkOpt struct {
	PeerID  int
	Message string
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

func (u *VkUsecase) Execute(
	ctx context.Context,
	opt VkOpt,
) (VkRes, error) {
	log.Printf("Received message from peer %d with text: '%s'", opt.PeerID, opt.Message)

	messageParams, err := u.parser.Parse(opt.Message)
	if err != nil {
		return VkRes{}, fmt.Errorf("Failed to parse message: %w", err)
	}

	if len(messageParams.Message) == 0 {
		return VkRes{}, fmt.Errorf("Message is empty")
	}

	img, err := u.t2i.Convert(ctx, messageParams.Message, messageParams.ImageParams)
	if err != nil {
		return VkRes{}, fmt.Errorf("Failed to convert message to image: %w", err)
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return VkRes{}, fmt.Errorf("Failed to encode image: %v", err)
	}

	resp, err := u.vk.UploadMessagesPhoto(opt.PeerID, &buf)
	if err != nil {
		return VkRes{}, fmt.Errorf("Failed to upload image: %w", err)
	}
	if len(resp) != 1 {
		return VkRes{}, fmt.Errorf("Invalid response length from UploadMessagesPhoto: %d", len(resp))
	}

	b := params.NewMessagesSendBuilder()
	b.RandomID(0)
	b.PeerID(opt.PeerID)
	b.Attachment(resp[0].ToAttachment())

	_, err = u.vk.MessagesSend(b.Params)
	if err != nil {
		return VkRes{}, fmt.Errorf("Failed to send message: %w", err)
	}

	return VkRes{}, nil
}
