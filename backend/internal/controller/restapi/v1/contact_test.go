package v1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/response"
	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

func TestSendContactMessage(t *testing.T) {
	contact := &contactUseCaseStub{}
	controller := NewController(nil, nil, nil, nil, contact, response.NewArtworkMapper("/media"))
	request := httptest.NewRequest(http.MethodPost, "/contact", strings.NewReader(`{"senderEmail":"sender@example.com","message":"Hello"}`))
	result := httptest.NewRecorder()

	controller.sendContactMessage(result, request)

	if result.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %q", result.Code, result.Body.String())
	}
	if contact.message.SenderEmail != "sender@example.com" || contact.message.Message != "Hello" {
		t.Fatalf("contact message = %#v", contact.message)
	}
}

type contactUseCaseStub struct {
	message entity.ContactMessage
}

func (u *contactUseCaseStub) Send(_ context.Context, message entity.ContactMessage) error {
	u.message = message
	return nil
}
