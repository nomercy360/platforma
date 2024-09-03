package payment

import (
	"context"
	"github.com/plutov/paypal/v4"
	"log"
)

type PaypalClient struct {
	client *paypal.Client
}

type PayPalRequest struct {
	PurchaseUnits      []paypal.PurchaseUnitRequest
	PaymentSource      *paypal.PaymentSource
	ApplicationContext *paypal.ApplicationContext
	Payer              *paypal.Payer
}

func NewPaypalClient(clientID, secret string, live bool) (*PaypalClient, error) {
	url := paypal.APIBaseLive
	if !live {
		url = paypal.APIBaseSandBox
	}
	c, err := paypal.NewClient(clientID, secret, url)

	if err != nil {
		return nil, err
	}

	return &PaypalClient{client: c}, nil
}

func (pc PaypalClient) CreatePaypalOrder(req PayPalRequest) (*paypal.Order, error) {
	ctx := context.Background()
	order, err := pc.client.CreateOrder(ctx, paypal.OrderIntentCapture, req.PurchaseUnits, req.PaymentSource, req.ApplicationContext)
	if err != nil {
		log.Printf("failed to create paypal order: %v", err)
		return nil, err
	}
	log.Printf("paypal order created: %v", order)
	return order, nil
}

func (pc PaypalClient) CapturePaypalOrder(orderID string) (*paypal.CaptureOrderResponse, error) {
	capture, err := pc.client.CaptureOrder(context.Background(), orderID, paypal.CaptureOrderRequest{})
	if err != nil {
		log.Printf("failed to capture paypal order: %v", err)
		return nil, err
	}

	log.Printf("paypal order captured: %v", capture)

	return capture, nil
}
