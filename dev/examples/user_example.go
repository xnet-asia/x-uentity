package main

import (
	"fmt"

	"github.com/xnetltd/x-uentity/factories"
	"github.com/xnetltd/x-uentity/handlers"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	factory := factories.NewFactory[User]().WithMiddleware(
		&handlers.ValidationMiddleware[User]{},
		handlers.NewLoggingMiddleware[User](func(message string) {
			fmt.Printf("MIDDLEWARE %s\n", message)
		}),
	)

	if _, err := factory.GetAuthHandler().Register("client-1", "token-123"); err != nil {
		panic(err)
	}

	server := factory.GetP2PServer()
	user := User{ID: "user-1", Name: "John Doe", Email: "john@example.com"}

	created, err := server.HandleHTTPFallback("token-123", &handlers.EntityRequest[User]{
		Action: "create",
		ID:     user.ID,
		Data:   user,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("AUTH CREATE code=%d success=%t user=%s\n", created.Code, created.Success, created.Single.Name)

	found, err := server.HandleHTTPFallback("", &handlers.EntityRequest[User]{
		Action: "get",
		ID:     user.ID,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("ANON GET   code=%d success=%t user=%s\n", found.Code, found.Success, found.Single.Name)

	denied, err := server.HandleHTTPFallback("", &handlers.EntityRequest[User]{
		Action: "delete",
		ID:     user.ID,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("ANON DELETE code=%d success=%t error=%q\n", denied.Code, denied.Success, denied.Error)
}
