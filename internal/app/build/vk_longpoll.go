package build

import (
	"context"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"

	"log"
	"os"

	"github.com/SevereCloud/vksdk/v3/api"
	"github.com/SevereCloud/vksdk/v3/events"
	"github.com/SevereCloud/vksdk/v3/longpoll-bot"
	"github.com/ashurbekovz/vktexbot/internal/app/parsers"
	"github.com/ashurbekovz/vktexbot/internal/app/usecases"
	"github.com/ashurbekovz/vktexbot/internal/pkg/latex2img"
	"github.com/ashurbekovz/vktexbot/internal/pkg/template2img"
	"github.com/shopspring/decimal"

	"os/signal"
	"syscall"
)

type VkApp struct {
	configPath string
	secretPath string
}

type Secret struct {
	GroupID string `yaml:"vk_group_id"`
	Token   string `yaml:"vk_token"`
}

type Config struct {
	Packages string `yaml:"packages"`
}

func NewVkApp(configPath string, secretPath string) *VkApp {
	return &VkApp{
		configPath: configPath,
		secretPath: secretPath,
	}
}

func (a *VkApp) Run() {
	config, err := load[Config](a.configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	secret, err := load[Secret](a.secretPath)
	if err != nil {
		log.Fatalf("Error loading secret: %v", err)
	}

	vk := api.NewVK(secret.Token)

	group, err := vk.GroupsGetByID(
		map[string]any{"group_id": secret.GroupID},
	)
	if err != nil {
		log.Fatalf("Error getting group info: %v", err)
	} else if len(group.Groups) != 1 {
		log.Fatalf("Unexpected number of groups found: %d", len(group.Groups))
	}

	lp, err := longpoll.NewLongPoll(vk, group.Groups[0].ID)
	if err != nil {
		log.Fatalf("Error creating Long Poll: %v", err)
	}

	l2i := latex2img.NewLatexToImgConverter("/tmp/", false, decimal.NewFromInt(400))
	t2i := template2img.NewLatexTemplateToImgConverter(&l2i, config.Packages)
	parser := parsers.NewVk(group.Groups[0].ID, group.Groups[0].ScreenName)
	vkUC := usecases.NewVkUsecase(vk, &t2i, parser)

	lp.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		opt := usecases.VkOpt{
			Message: obj.Message.Text,
			PeerID:  obj.Message.PeerID,
			Payload: obj.Message.Payload,
		}

		_, err := vkUC.Execute(ctx, opt)
		if err != nil {
			log.Printf("Error during handling MessageNew event: %v", err)
			return
		}
	})

	lp.MessageEdit(func(ctx context.Context, obj events.MessageEditObject) {
		opt := usecases.VkOpt{
			Message: obj.Text,
			PeerID:  obj.PeerID,
			Payload: obj.Payload,
		}

		_, err := vkUC.Execute(ctx, opt)
		if err != nil {
			log.Printf("Error during handling MessageEdit event: %v", err)
			return
		}
	})

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := <-shutdown
		log.Printf("received shutdown signal: %v", sig)
		cancel()
	}()

	log.Println("Start VK Long Poll Server")
	if err := lp.RunWithContext(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("VK Long Poll Server: graceful shutdown")
		} else {
			log.Fatalf("VK Long Poll Server crashed: %v", err)
		}
	}

	log.Println("VK Long Poll Server stopped")
}

func load[config any](path string) (*config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error opening config file: %w", err)
	}
	defer file.Close()

	var cfg config
	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("Error decoding config file: %w", err)
	}

	return &cfg, nil
}
