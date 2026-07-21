package handlers

import (
	"reflect"
	"testing"

	"github.com/xnetltd/x-uentity/repositories"
)

type handlerTestEntity struct {
	Name string `json:"name"`
}

type recordingMiddleware struct {
	name  string
	calls *[]string
}

func (m recordingMiddleware) Ingress(_ *ClientAuth, _ *EntityRequest[handlerTestEntity]) error {
	*m.calls = append(*m.calls, "in:"+m.name)
	return nil
}

func (m recordingMiddleware) Egress(_ *ClientAuth, _ *EntityResponse[handlerTestEntity]) error {
	*m.calls = append(*m.calls, "out:"+m.name)
	return nil
}

func TestEntityHandlerPermissionsAndCRUD(t *testing.T) {
	repo := repositories.NewInMemoryRepository[handlerTestEntity]()
	handler := NewEntityHandler(repo, NewMiddlewareChain[handlerTestEntity]())
	entity := handlerTestEntity{Name: "Ada"}

	resp, err := handler.Handle(nil, &EntityRequest[handlerTestEntity]{Action: "create", ID: "user-1", Data: entity})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Code != 403 {
		t.Fatalf("anonymous create code = %d, want 403", resp.Code)
	}
	t.Logf("ANON CREATE code=%d success=%t error=%q", resp.Code, resp.Success, resp.Error)

	auth := &ClientAuth{ID: "client-1", IsAuth: true}
	resp, err = handler.Handle(auth, &EntityRequest[handlerTestEntity]{Action: "create", ID: "user-1", Data: entity})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Success || resp.Code != 201 {
		t.Fatalf("authenticated create response = %+v", resp)
	}
	t.Logf("AUTH CREATE code=%d success=%t entity=%+v", resp.Code, resp.Success, *resp.Single)

	resp, err = handler.Handle(nil, &EntityRequest[handlerTestEntity]{Action: "get", ID: "user-1"})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Success || resp.Single == nil || *resp.Single != entity {
		t.Fatalf("anonymous get response = %+v", resp)
	}
	t.Logf("ANON GET code=%d success=%t entity=%+v", resp.Code, resp.Success, *resp.Single)
}

func TestMiddlewareOrder(t *testing.T) {
	var calls []string
	first := recordingMiddleware{name: "first", calls: &calls}
	second := recordingMiddleware{name: "second", calls: &calls}
	handler := NewEntityHandler(
		repositories.NewInMemoryRepository[handlerTestEntity](),
		NewMiddlewareChain[handlerTestEntity](first, second),
	)

	if _, err := handler.Handle(nil, &EntityRequest[handlerTestEntity]{Action: "query"}); err != nil {
		t.Fatal(err)
	}

	want := []string{"in:first", "in:second", "out:second", "out:first"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("middleware calls = %v, want %v", calls, want)
	}
}
