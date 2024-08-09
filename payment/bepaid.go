package payment

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type BepaidAVSCVCVerification struct {
	Enabled bool `json:"enabled"`
}

type BepaidCardOnFile struct {
	Initiator string `json:"initiator,omitempty"`
	Type      string `json:"type,omitempty"`
}

type BepaidSaveCardToggle struct {
	Display bool   `json:"display,omitempty"`
	Text    string `json:"text,omitempty"`
	Hint    string `json:"hint,omitempty"`
}

type BepaidAnotherCardToggle struct {
	Display bool `json:"display,omitempty"`
}

type BepaidSettings struct {
	ReturnUrl           string                  `json:"return_url,omitempty"`
	SuccessUrl          string                  `json:"success_url,omitempty"`
	DeclineUrl          string                  `json:"decline_url,omitempty"`
	FailUrl             string                  `json:"fail_url,omitempty"`
	CancelUrl           string                  `json:"cancel_url,omitempty"`
	NotificationUrl     string                  `json:"notification_url,omitempty"`
	ButtonText          string                  `json:"button_text,omitempty"`
	ButtonNextText      string                  `json:"button_next_text,omitempty"`
	Language            string                  `json:"language,omitempty"`
	CardNotificationUrl string                  `json:"card_notification_url,omitempty"`
	AutoPay             bool                    `json:"auto_pay,omitempty,omitempty"`
	SaveCardToggle      BepaidSaveCardToggle    `json:"save_card_toggle,omitempty"`
	AnotherCardToggle   BepaidAnotherCardToggle `json:"another_card_toggle,omitempty"`
	PaymentMethod       BepaidPaymentMethod     `json:"payment_method,omitempty"`
	AutoReturn          string                  `json:"auto_return,omitempty"`
	WidgetStyle         interface{}             `json:"style,omitempty"`
}

type BepaidPaymentMethod struct {
	Types []string `json:"types,omitempty"`
}

type BepaidCustomer struct {
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Address   string `json:"address,omitempty"`
	City      string `json:"city,omitempty"`
	State     string `json:"state,omitempty"`
	ZIP       string `json:"zip,omitempty"`
	Country   string `json:"country,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type BepaidOrder struct {
	Currency       string         `json:"currency"`
	Amount         int            `json:"amount"`
	Description    string         `json:"description"`
	AdditionalData AdditionalData `json:"additional_data"`
	TrackingID     string         `json:"tracking_id"`
}

type AdditionalData struct {
	ReceiptText        []string                 `json:"receipt_text,omitempty"`
	Contract           []string                 `json:"contract,omitempty"`
	AVSCVCVerification BepaidAVSCVCVerification `json:"avs_cvc_verification,omitempty"`
}

type BepaidCheckout struct {
	TransactionType string           `json:"transaction_type"`
	Test            bool             `json:"test,omitempty"`
	Iframe          bool             `json:"iframe,omitempty"`
	CardOnFile      BepaidCardOnFile `json:"card_on_file,omitempty"`
	Settings        BepaidSettings   `json:"settings,omitempty"`
	Customer        BepaidCustomer   `json:"customer,omitempty"`
	Attempts        int              `json:"attempts,omitempty"`
	Order           BepaidOrder      `json:"order"`
}

type BepaidTokenRequest struct {
	Checkout BepaidCheckout `json:"checkout"`
}

type BepaidTokenResponse struct {
	Checkout BepaidResponseCheckout `json:"checkout"`
}

type BepaidResponseCheckout struct {
	Token       string `json:"token"`
	RedirectUrl string `json:"redirect_url"`
}

type BepaidNotification struct {
	Transaction BepaidTransaction `json:"transaction"`
}

type BepaidTransaction struct {
	Uid                 string     `json:"uid"`
	Status              string     `json:"status"`
	Amount              int        `json:"amount"`
	Currency            string     `json:"currency"`
	Description         string     `json:"description"`
	Type                string     `json:"type"`
	PaymentMethodType   string     `json:"payment_method_type"`
	TrackingId          string     `json:"tracking_id"`
	Message             string     `json:"message"`
	Test                bool       `json:"test"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	PaidAt              time.Time  `json:"paid_at"`
	ExpiredAt           *time.Time `json:"expired_at"`
	RecurringType       *time.Time `json:"recurring_type"`
	ClosedAt            *time.Time `json:"closed_at"`
	SettledAt           *time.Time `json:"settled_at"`
	ManuallyCorrectedAt *time.Time `json:"manually_corrected_at"`
	Language            string     `json:"language"`
	ID                  string     `json:"id"`
}

func CreatePaymentToken(request BepaidTokenRequest, apiURL, shopID, shopSecret string) (*BepaidTokenResponse, error) {
	auth := base64.StdEncoding.EncodeToString([]byte(shopID + ":" + shopSecret))

	client := &http.Client{}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Version", "2")
	req.Header.Set("Authorization", "Basic "+auth)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New("failed to create payment token")
	}

	var tokenResp BepaidTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}
