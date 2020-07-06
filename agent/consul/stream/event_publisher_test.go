package stream

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testTopic Topic = 999

func TestEventPublisher_PublishChangesAndSubscribe_WithSnapshot(t *testing.T) {
	subscription := &SubscribeRequest{
		Topic: testTopic,
		Key:   "sub-key",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	publisher := NewEventPublisher(ctx, newTestSnapshotHandlers(), 0)
	sub, err := publisher.Subscribe(ctx, subscription)
	require.NoError(t, err)
	eventCh := consumeSubscription(sub)

	result := nextResult(t, eventCh)
	require.NoError(t, result.Err)
	expected := []Event{{Payload: "snapshot-event-payload", Key: "sub-key"}}
	require.Equal(t, expected, result.Events)

	result = nextResult(t, eventCh)
	require.Len(t, result.Events, 1)
	require.True(t, result.Events[0].IsEndOfSnapshot())

	// Now subscriber should block waiting for updates
	assertNoResult(t, eventCh)

	events := []Event{{
		Topic:   testTopic,
		Key:     "sub-key",
		Payload: "the-published-event-payload",
	}}
	publisher.PublishEvents(events)

	// Subscriber should see the published event
	result = nextResult(t, eventCh)
	require.NoError(t, result.Err)
	expected = []Event{{Payload: "the-published-event-payload", Key: "sub-key", Topic: testTopic}}
	require.Equal(t, expected, result.Events)
}

func newTestSnapshotHandlers() SnapshotHandlers {
	return SnapshotHandlers{
		testTopic: func(req *SubscribeRequest, buf SnapshotAppender) (uint64, error) {
			if req.Topic != testTopic {
				return 0, fmt.Errorf("unexpected topic: %v", req.Topic)
			}
			buf.Append([]Event{{Payload: "snapshot-event-payload", Key: "sub-key"}})
			return 1, nil
		},
	}
}

func consumeSubscription(sub *Subscription) <-chan subNextResult {
	eventCh := make(chan subNextResult, 1)
	go func() {
		for {
			es, err := sub.Next()
			eventCh <- subNextResult{
				Events: es,
				Err:    err,
			}
			if err != nil {
				return
			}
		}
	}()
	return eventCh
}

type subNextResult struct {
	Events []Event
	Err    error
}

func nextResult(t *testing.T, eventCh <-chan subNextResult) subNextResult {
	t.Helper()
	select {
	case next := <-eventCh:
		return next
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("no event after 100ms")
	}
	return subNextResult{}
}

func assertNoResult(t *testing.T, eventCh <-chan subNextResult) {
	t.Helper()
	select {
	case next := <-eventCh:
		require.NoError(t, next.Err)
		require.Len(t, next.Events, 1)
		t.Fatalf("received unexpected event: %#v", next.Events[0].Payload)
	case <-time.After(100 * time.Millisecond):
	}
}
