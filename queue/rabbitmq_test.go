package queue

import (
	"context"
	"fmt"
	"testing"
)

func TestRabbitMQ(t *testing.T) {
	t.Skip("skip")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conf := &RabbitMQConfig{
		Addr:        "localhost",
		Port:        5672,
		VirtualHost: "/",
		User:        "guest",
		Password:    "guest",
		QueueName:   "test_queue_123456",
	}
	queue, err := NewRabbitMQ(ctx, conf)
	if err != nil {
		t.Fatalf("Error NewRabbitMQ: %v", err)
	}
	fmt.Println("queue: ", queue.queue.Name)
	defer queue.Close()
	defer queue.Delete()

	options := PublishOptions{}
	options.Priority = 1
	err = queue.ProduceWithCtx(ctx, "hello world", options)
	if err != nil {
		t.Fatalf("Error producing to Rabbit queue: %v", err)
	}
	options.Priority = 2
	err = queue.ProduceWithCtx(ctx, "hi there", options)
	if err != nil {
		t.Fatalf("Error producing to Rabbit queue: %v", err)
	}
	options.Priority = 2
	err = queue.ProduceWithCtx(ctx, "it's over", options)
	if err != nil {
		t.Fatalf("Error producing to Rabbit queue: %v", err)
	}

	data, err := queue.ConsumeWithCtx(ctx)
	if err != nil {
		t.Fatalf("Error consumeing from Rabbit queue: %v", err)
		return
	}
	if data != "hi there" {
		t.Fatalf("Expected message 'hi there' from Rabbit queue, but got '%s'", data)
	}

	rs := make([]string, 0)
	// Consume a message from Redis queue
	err = queue.ConsumeFuncWithCtx(ctx, func(data string) bool {
		t.Log(data)
		rs = append(rs, data)
		if len(rs) >= 2 {
			cancel()
		}
		return true
	})
	if err != nil && err != ctx.Err() {
		t.Fatalf("Error consuming from Rabbit queue: %v", err)
	}
	if len(rs) != 2 {
		t.Fatalf("Expected message num 2 from Rabbit queue, but got '%d'", len(rs))
	}
	if rs[0] != "it's over" || rs[1] != "hello world" {
		t.Fatalf("Expected message list 'it's over,hello world' from Rabbit queue, but got '%s,%s'", rs[0], rs[1])
	}
}
