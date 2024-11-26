package graph

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 메시지를 저장할 배열과 Mutex를 선언하여 동시성 문제를 방지
var messages []*string
var mutex sync.RWMutex

// UpdateMessage is the resolver for the updateMessage field.
func (r *mutationResolver) UpdateMessage(_ context.Context, input string) (string, error) {
	mutex.Lock()
	messages = append(messages, &input)
	mutex.Unlock()
	return "Message added: " + input, nil
}

// DeleteMessage is the resolver for the deleteMessage field.
func (r *mutationResolver) DeleteMessage(_ context.Context, index int) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	// Index 범위 검사
	if index < 0 || index >= len(messages) {
		return "", fmt.Errorf("invalid index: %d", index)
	}

	deletedMessage := messages[index]
	messages = append(messages[:index], messages[index+1:]...) // 메시지 삭제
	return "Deleted message: " + *deletedMessage, nil
}

// Hello is the resolver for the hello field.
func (r *queryResolver) Hello(_ context.Context) (string, error) {
	return "Hello, gqlgen!", nil
}

// GetMessages is the resolver for the getMessages field.
func (r *queryResolver) GetMessages(_ context.Context) ([]*string, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	return messages, nil
}

// MessageStream is the resolver for the messageStream field.
func (r *subscriptionResolver) MessageStream(_ context.Context) (<-chan string, error) {
	messageChan := make(chan string)
	go func() {
		defer close(messageChan)
		for i := 0; i < 5; i++ {
			messageChan <- "Message " + time.Now().String()
			time.Sleep(1 * time.Second)
		}
	}()
	return messageChan, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
