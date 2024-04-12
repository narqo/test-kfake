package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kfake"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	testTopic = "test-kgo-topic"

	testTopicPartition = 1
)

func TestListOffsetsAfterMilli(t *testing.T) {
	t.Parallel()

	t.Run("kfake", func(t *testing.T) {
		cluster := kfake.MustCluster(
			kfake.SeedTopics(-1, testTopic),
		)
		t.Cleanup(cluster.Close)

		client, err := kgo.NewClient(
			kgo.SeedBrokers(cluster.ListenAddrs()...),
			kgo.RecordPartitioner(kgo.ManualPartitioner()),
		)
		require.NoError(t, err)
		t.Cleanup(client.Close)

		testListOffsetsAfterMilli(t, client)
	})

	t.Run("kafka", func(t *testing.T) {
		seedBroker := "127.0.0.1:29092"

		client, err := kgo.NewClient(
			kgo.SeedBrokers(seedBroker),
			kgo.RecordPartitioner(kgo.ManualPartitioner()),
			kgo.AllowAutoTopicCreation(),
		)
		require.NoError(t, err)
		t.Cleanup(client.Close)

		// remove the topic on start
		kadm.NewClient(client).DeleteTopic(context.Background(), testTopic)

		testListOffsetsAfterMilli(t, client)
	})
}

func testListOffsetsAfterMilli(t *testing.T, client *kgo.Client) {
	ctx := context.Background()

	var lastRec *kgo.Record
	for i := range 5 {
		if i > 0 {
			time.Sleep(100 * time.Millisecond)
		}

		rec := kgo.StringRecord(fmt.Sprintf("rec-%d", i))
		rec.Topic = testTopic
		rec.Partition = testTopicPartition

		var err error
		lastRec, err = client.ProduceSync(ctx, rec).First()
		require.NoError(t, err)
	}

	require.NotNil(t, lastRec)

	admClient := kadm.NewClient(client)

	// list offsets AT last record's timestamp
	listed, err := admClient.ListOffsetsAfterMilli(ctx, lastRec.Timestamp.UnixMilli(), testTopic)
	require.NoError(t, err)

	offset, ok := listed.Lookup(testTopic, testTopicPartition)
	require.True(t, ok)
	if !assert.EqualValues(t, 4, offset.Offset) {
		t.Logf("offset %+v", offset)
	}

	// list offsets AFTER last record's timestamp
	listed, err = admClient.ListOffsetsAfterMilli(ctx, lastRec.Timestamp.Add(time.Millisecond).UnixMilli(), testTopic)
	require.NoError(t, err)

	offset, ok = listed.Lookup(testTopic, testTopicPartition)
	require.True(t, ok)
	if !assert.EqualValues(t, 5, offset.Offset) {
		t.Logf("offset %+v", offset)
	}
}
