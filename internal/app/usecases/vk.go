package usecases

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"log"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/api/params"
	"github.com/ashurbekovz/vktexbot/internal/pkg/template2img"
)

type VkOpt struct {
	PeerID  int
	Message string
}

type VkRes struct {
}

type VkUsecase struct {
	vk  *api.VK
	t2i *template2img.LatexTemplateToImgConverter
}

func NewVkUsecase(
	vk *api.VK,
	t2i *template2img.LatexTemplateToImgConverter,
) *VkUsecase {
	return &VkUsecase{
		vk:  vk,
		t2i: t2i,
	}
}

func (u *VkUsecase) Execute(
	ctx context.Context,
	opt VkOpt,
) (VkRes, error) {
	log.Printf("Received message from peer %d with text: '%s'", opt.PeerID, opt.Message)

	imgParams, err := template2img.NewImageParams()
	if err != nil {
		return VkRes{}, fmt.Errorf("Failed to create image parameters: %w", err)
	}

	img, err := u.t2i.Convert(ctx, opt.Message, imgParams)
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
