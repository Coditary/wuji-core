package core_test

import (
	"context"
	"strings"
	"testing"

	"github.com/coditary/wuji-core/pkg/config"
	"github.com/coditary/wuji-core/pkg/core"
	"github.com/coditary/wuji-core/pkg/driver"
	"github.com/coditary/wuji-core/pkg/driver/dummy"
)

func testCore(t *testing.T) *core.Core {
	t.Helper()
	c, err := core.New(core.Config{
		AppConfig:       &config.Config{DefaultDriver: dummy.DriverID},
		DefaultDriverID: dummy.DriverID,
		Lazy:            true,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

func TestCoreGenerateText(t *testing.T) {
	c := testCore(t)
	defer c.Close()

	resp, err := c.GenerateText(context.Background(), "", driver.TextRequest{
		Prompt: "test prompt", MaxTokens: 128, Temperature: 0.5,
	})
	if err != nil {
		t.Fatalf("GenerateText: %v", err)
	}
	if resp.Text == "" {
		t.Fatal("expected non-empty text response")
	}
}

func TestCoreGenerateTextWithMessages(t *testing.T) {
	c := testCore(t)
	defer c.Close()

	resp, err := c.GenerateText(context.Background(), "", driver.TextRequest{
		Messages: []driver.ChatMessage{
			{Role: driver.ChatRoleSystem, Content: "be concise"},
			{Role: driver.ChatRoleUser, Content: "hello"},
		},
		MaxTokens: 128,
	})
	if err != nil {
		t.Fatalf("GenerateText: %v", err)
	}
	if resp.Text == "" {
		t.Fatal("expected non-empty text response")
	}
	if !strings.Contains(resp.Text, "[dummy:messages]") {
		t.Fatalf("expected messages in dummy response, got %q", resp.Text)
	}
}

func TestCoreGenerateTextFromImage(t *testing.T) {
	c := testCore(t)
	defer c.Close()

	resp, err := c.GenerateText(context.Background(), "", driver.TextRequest{
		InputMode: driver.TextInputImage,
		MediaPath: "photo.jpg",
	})
	if err != nil {
		t.Fatalf("GenerateText: %v", err)
	}
	if resp.Text == "" {
		t.Fatal("expected non-empty text response")
	}
}

func TestCoreGenerateImage(t *testing.T) {
	c := testCore(t)
	defer c.Close()

	resp, err := c.GenerateImage(context.Background(), "", driver.ImageRequest{
		Prompt: "a cat", Width: 512, Height: 512,
	})
	if err != nil {
		t.Fatalf("GenerateImage: %v", err)
	}
	if resp.Path == "" {
		t.Fatal("expected non-empty image path")
	}
}

func TestCoreListDrivers(t *testing.T) {
	c := testCore(t)
	defer c.Close()

	drivers := c.ListDrivers()
	if len(drivers) < 1 {
		t.Fatalf("expected at least 1 driver, got %d", len(drivers))
	}

	found := false
	for _, d := range drivers {
		if d.ID == dummy.DriverID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected dummy driver to be registered")
	}
}

func TestCoreLazyAllowsMissingDefaultDriver(t *testing.T) {
	c, err := core.New(core.Config{
		AppConfig: &config.Config{DefaultDriver: "echo"},
		Lazy:      true,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer c.Close()

	if c.DefaultDriverID() != "echo" {
		t.Fatalf("expected default driver echo, got %q", c.DefaultDriverID())
	}
}

func TestCoreManageDataset(t *testing.T) {
	c := testCore(t)
	defer c.Close()

	resp, err := c.ManageDataset(context.Background(), "", driver.DatasetRequest{
		Action: driver.DatasetList,
	})
	if err != nil {
		t.Fatalf("ManageDataset: %v", err)
	}
	if len(resp.Datasets) == 0 {
		t.Fatal("expected at least one dataset")
	}
}
