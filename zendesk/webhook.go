package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Webhook is struct for webhook payload.
// https://developer.zendesk.com/api-reference/event-connectors/webhooks/webhooks/
type Webhook struct {
	Authentication *WebhookAuthentication `json:"authentication,omitempty"`
	CreatedAt      time.Time              `json:"created_at,omitempty"`
	CreatedBy      string                 `json:"created_by,omitempty"`
	Description    string                 `json:"description,omitempty"`
	Endpoint       string                 `json:"endpoint"`
	ExternalSource interface{}            `json:"external_source,omitempty"`
	HTTPMethod     string                 `json:"http_method"`
	ID             string                 `json:"id,omitempty"`
	Name           string                 `json:"name"`
	RequestFormat  string                 `json:"request_format"`
	SigningSecret  *WebhookSigningSecret  `json:"signing_secret,omitempty"`
	Status         string                 `json:"status"`
	Subscriptions  []string               `json:"subscriptions,omitempty"`
	UpdatedAt      time.Time              `json:"updated_at,omitempty"`
	UpdatedBy      string                 `json:"updated_by,omitempty"`
}

type WebhookAuthentication struct {
	Type        string      `json:"type"`
	Data        interface{} `json:"data"`
	AddPosition string      `json:"add_position"`
}

type WebhookSigningSecret struct {
	Algorithm string `json:"algorithm"`
	Secret    string `json:"secret"`
}

// WebhookListOptions is options for ListWebhooks
//
// ref: https://developer.zendesk.com/api-reference/webhooks/webhooks-api/webhooks/#list-webhooks
type WebhookListOptions struct {
	PageOptions
	FilterNameContains string `url:"filter[name_contains],omitempty"`
	FilterStatus       string `url:"filter[status],omitempty"`
	PageAfter          string `url:"page[after],omitempty"`
	PageBefore         string `url:"page[before],omitempty"`
	PageSize           string `url:"page[size],omitempty"`
	Sort               string `url:"sort,omitempty"`
}

type WebhookAPI interface {
	ListWebhooks(ctx context.Context, opts *WebhookListOptions) ([]Webhook, Page, error)
	CreateWebhook(ctx context.Context, hook *Webhook) (*Webhook, error)
	GetWebhook(ctx context.Context, webhookID string) (*Webhook, error)
	UpdateWebhook(ctx context.Context, webhookID string, hook *Webhook) error
	DeleteWebhook(ctx context.Context, webhookID string) error
	GetWebhookSigningSecret(ctx context.Context, webhookID string) (*WebhookSigningSecret, error)
}

// ListWebhooks lists webhooks.
//
// https://developer.zendesk.com/api-reference/webhooks/webhooks-api/webhooks/#list-webhooks
func (z *Client) ListWebhooks(ctx context.Context, opts *WebhookListOptions) ([]Webhook, Page, error) {
	var data struct {
		Webhooks []Webhook `json:"webhooks"`
		Page
	}

	tmp := opts
	if tmp == nil {
		tmp = &WebhookListOptions{}
	}

	u, err := addOptions("/webhooks", tmp)
	if err != nil {
		return nil, Page{}, err
	}

	body, err := z.get(ctx, u)
	if err != nil {
		return nil, Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, Page{}, err
	}
	return data.Webhooks, data.Page, nil
}

// CreateWebhook creates new webhook.
//
// https://developer.zendesk.com/api-reference/event-connectors/webhooks/webhooks/#create-or-clone-webhook
func (z *Client) CreateWebhook(ctx context.Context, hook *Webhook) (*Webhook, error) {
	var data, result struct {
		Webhook *Webhook `json:"webhook"`
	}
	data.Webhook = hook

	body, err := z.post(ctx, "/webhooks", data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Webhook, nil
}

// GetWebhook gets a specified webhook.
//
// https://developer.zendesk.com/api-reference/event-connectors/webhooks/webhooks/#show-webhook
func (z *Client) GetWebhook(ctx context.Context, webhookID string) (*Webhook, error) {
	var result struct {
		Webhook *Webhook `json:"webhook"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/webhooks/%s", webhookID))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Webhook, nil
}

// UpdateWebhook updates a webhook with the specified webhook.
//
// https://developer.zendesk.com/api-reference/event-connectors/webhooks/webhooks/#update-webhook
func (z *Client) UpdateWebhook(ctx context.Context, webhookID string, hook *Webhook) error {
	var data struct {
		Webhook *Webhook `json:"webhook"`
	}
	data.Webhook = hook

	_, err := z.put(ctx, fmt.Sprintf("/webhooks/%s", webhookID), data)
	if err != nil {
		return err
	}

	return nil
}

// DeleteWebhook deletes the specified webhook.
//
// https://developer.zendesk.com/api-reference/event-connectors/webhooks/webhooks/#delete-webhook
func (z *Client) DeleteWebhook(ctx context.Context, webhookID string) error {
	err := z.delete(ctx, fmt.Sprintf("/webhooks/%s", webhookID))
	if err != nil {
		return err
	}

	return nil
}

// GetWebhookSigningSecret gets the signing secret of specified webhook.
//
// https://developer.zendesk.com/api-reference/event-connectors/webhooks/webhooks/#show-webhook-signing-secret
func (z *Client) GetWebhookSigningSecret(ctx context.Context, webhookID string) (*WebhookSigningSecret, error) {
	var result struct {
		SigningSecret *WebhookSigningSecret `json:"signing_secret"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/webhooks/%s/signing_secret", webhookID))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.SigningSecret, nil
}
