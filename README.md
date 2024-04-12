# test-kfake

An example, to illustrate the https://github.com/twmb/franz-go/issues/706

## Usage

1\. Start local Kafka cluster


```
% docker compose up kafka
```

2\. Run Go tests

```
% go test -v ./
=== RUN   TestListOffsetsAfterMilli
=== RUN   TestListOffsetsAfterMilli/kfake
    main_test.go:87:
        	Error Trace:	/Users/v/Documents/Code/test-kfake/main_test.go:87
        	            				/Users/v/Documents/Code/test-kfake/main_test.go:38
        	Error:      	Not equal:
        	            	expected: int(4)
        	            	actual  : int64(0)
        	Test:       	TestListOffsetsAfterMilli/kfake
    main_test.go:88: offset {Topic:test-kgo-topic Partition:1 Timestamp:-1 Offset:0 LeaderEpoch:0 Err:<nil>}
    main_test.go:97:
        	Error Trace:	/Users/v/Documents/Code/test-kfake/main_test.go:97
        	            				/Users/v/Documents/Code/test-kfake/main_test.go:38
        	Error:      	Not equal:
        	            	expected: int(5)
        	            	actual  : int64(0)
        	Test:       	TestListOffsetsAfterMilli/kfake
    main_test.go:98: offset {Topic:test-kgo-topic Partition:1 Timestamp:-1 Offset:0 LeaderEpoch:0 Err:<nil>}
=== RUN   TestListOffsetsAfterMilli/kafka
--- FAIL: TestListOffsetsAfterMilli (1.16s)
    --- FAIL: TestListOffsetsAfterMilli/kfake (0.41s)
    --- PASS: TestListOffsetsAfterMilli/kafka (0.75s)
FAIL
FAIL	github.com/narqo/test-kfake	1.301s
FAIL
```

The results above show that the same test passed against Kafka, but failed when run against a `kfake` cluster.
