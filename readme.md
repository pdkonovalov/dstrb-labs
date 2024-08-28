# Distributed Systems labs.

The [course](http://nil.csail.mit.edu/6.5840/2023/index.html "MIT 6.5840: Distributed Systems") provide framework and tests to simulate real distributed systems on a single Linux computer. Different computers is different programs, that communicate via RPC over Unix sockets.

## Lab 1: MapReduce

Simple example, how it is work. We have some books in *src/main* directory. Lets search all rows in this files, that contains "map" or "reduce" word. We can use distributed grep for it.

```
cd src/main
go build -buildmode=plugin ../mrapps/grep.go
go run mrcoordinator.go pg-*.txt
```

Now coordinator wait workers. We can run some count of workers in different terminal windows. 

```
cd src/main
go run mrworker.go grep.so
```

Regardless of workers behavior, coordinator will control, that MapReduce process will be finished correctly and we can see output.

```
cat mr-out-*
now and then directed me, and I possessed a map of the country; but I pg-frankenstein.txt:4368
shock, but entirely reduced to thin ribbons of wood.  I never beheld pg-frankenstein.txt:939
...
```

Also we can run the labs test script:

```
bash test-mr.sh
*** Starting wc test.
--- wc test: PASS
*** Starting indexer test.
--- indexer test: PASS
*** Starting map parallelism test.
--- map parallelism test: PASS
*** Starting reduce parallelism test.
--- reduce parallelism test: PASS
*** Starting job count test.
--- job count test: PASS
*** Starting early exit test.
--- early exit test: PASS
*** Starting crash test.
--- crash test: PASS
*** PASSED ALL TESTS
```

More information about MapReduce can be found [here](https://static.googleusercontent.com/media/research.google.com/ru//archive/mapreduce-osdi04.pdf "MapReduce: Simplified Data Processing on Large Clusters").

## Lab 2: Raft

In this lab we implement raft package, following the design in this [paper](http://nil.csail.mit.edu/6.5840/2023/papers/raft-extended.pdf "In Search of an Understandable Consensus Algorithm").  If we have cluster of servers, we can use raft to replicate sequence of requests.

```
rf := raft.Make(peers, me, persister, applyCh)
```

This function create raft object. Peers is addreses of servers in cluster. Me is index of this server. Persister is object for save data before crash. In apply channel raft will send replicated commands.

```
index, term, isLeader := rf.Start(command)
```

If this server is leader and larger part of cluster is availible we will see the command in apply chanel. [This](https://thesecretlivesofdata.com/raft/ "thesecretlivesofdata raft") is graphic visualisation how it is work.

Lab tests:

```
cd src/raft
go test
Test (2A): initial election ...
  ... Passed --   3.1  3   60   16328    0
Test (2A): election after network failure ...
  ... Passed --   4.6  3  150   27089    0
Test (2A): multiple elections ...
  ... Passed --   5.6  7  726  125561    0
Test (2B): basic agreement ...
  ... Passed --   0.5  3   16    4338    3
Test (2B): RPC byte count ...
  ... Passed --   1.4  3   48  113758   11
Test (2B): test progressive failure of followers ...
  ... Passed --   4.6  3  142   26895    3
Test (2B): test failure of leaders ...
  ... Passed --   4.8  3  194   39524    3
Test (2B): agreement after follower reconnects ...
  ... Passed --   5.5  3  130   32670    8
Test (2B): no agreement if too many followers disconnect ...
  ... Passed --   3.3  5  260   47962    3
Test (2B): concurrent Start()s ...
  ... Passed --   0.7  3   24    6885    6
Test (2B): rejoin of partitioned leader ...
  ... Passed --   3.7  3  142   31757    4
Test (2B): leader backs up quickly over incorrect follower logs ...
  ... Passed --  16.9  5 2280 1568234  102
Test (2B): RPC counts aren't too high ...
  ... Passed --   2.1  3   64   20022   12
Test (2C): basic persistence ...
  ... Passed --   2.9  3   84   20635    6
Test (2C): more persistence ...
  ... Passed --  18.8  5 1388  255070   16
Test (2C): partitioned leader and one follower crash, leader restarts ...
  ... Passed --   1.5  3   40    9548    4
Test (2C): Figure 8 ...
  ... Passed --  31.0  5 1696  355799   66
Test (2C): unreliable agreement ...
  ... Passed --   1.3  5 1000  318091  246
Test (2C): Figure 8 (unreliable) ...
  ... Passed --  30.4  5 8792 15349122  325
Test (2C): churn ...
  ... Passed --  16.6  5 10252 32907300 2449
Test (2C): unreliable churn ...
  ... Passed --  16.2  5 4304 5617578  948
Test (2D): snapshots basic ...
  ... Passed --   4.0  3  514  197856  220
Test (2D): install snapshots (disconnect) ...
  ... Passed --  74.7  3 2270  914046  358
Test (2D): install snapshots (disconnect+unreliable) ...
  ... Passed --  79.9  3 2460  869651  301
Test (2D): install snapshots (crash) ...
  ... Passed --  26.9  3 1222  560555  326
Test (2D): install snapshots (unreliable+crash) ...
  ... Passed --  42.7  3 1520  707611  324
Test (2D): crash and restart all servers ...
  ... Passed --   8.3  3  314   85120   65
Test (2D): snapshot initialization after crash ...
  ... Passed --   2.1  3   80   20480   14
PASS
ok      6.5840/raft    414.146s
```

## Lab 3: Fault-tolerant Key/Value Service

In this lab we will use our raft package for build distributed kv-service, that need to be fault-tolerant and linearizable. Clients can send three requests: get, put and append. Keys and values is strings. Get - fetch current value or return empty string for non-existing key. Put - change value to new. Append - append value to exist or put new.

```
ck := kvserver.MakeClerk(servers)
```

This function make "clerk" object. Clerk manage client requests. Servers is addreses of all servers in cluster. Clerk will send rpc to this servers, while request wil not be executed.

```
ck.Put("distributed", "MapReduce")
ck.Append("distributed", " and Raft")
result := ck.Get("distributed")
```

If no other client iteracted with this key, we will see "MapReduce and Raft" and this value was succesfully replicated. 

But what happend if other clients send put or append to this key at the same time? Let another client send request too.

```
// client A
time.Now() // 12:01
ck.Put("distributed", "MapReduce")
time.Now() // 12:03
ck.Append("distributed", " and Raft")
time.Now() // 12:07
result := ck.Get("distributed")
-----------------------
// client B
time.Now() // 12:04
ck.Append("distributed", " and Kv-service")
time.Now() // 12:06
```

The result can be "MapReduce and Kv-service and Raft" or "MapReduce and Raft and Kv-service", but not another.

Our kv-service provide linearizability for one key. This means that in the servers side, history of first success servers responces to clerks rpc for one key - linear.

```
put "MapReduce"         -> ok                                           12:02
append "and Raft"       -> ok                                           12:04
append "and Kv-service" -> ok                                           12:05
get "distributed"       -> "MapReduce and Raft and Kv-service"          12:08
-----------------------
put "MapReduce"         -> ok                                           12:02
append "and Kv-service" -> ok                                           12:05
append "and Raft"       -> ok                                           12:06
get "distributed"       -> "MapReduce and Kv-service and Raft"          12:08
```

You can find more about this property in this [note](http://nil.csail.mit.edu/6.5840/2024/papers/linearizability-faq.txt "linearizability-faq") from the course authors.

Lab tests:

```
cd src/kvraft
go test
Test: one client (3A) ...
  ... Passed --  15.2  5 11102 2214
Test: ops complete fast enough (3A) ...
  ... Passed --   3.5  3  3014    0
Test: many clients (3A) ...
  ... Passed --  15.4  5 13246 2577
Test: unreliable net, many clients (3A) ...
  ... Passed --  16.5  5  6547 1067
Test: concurrent append to same key, unreliable (3A) ...
  ... Passed --   1.4  3   219   52
Test: progress in majority (3A) ...
  ... Passed --   0.3  5    41    2
Test: no progress in minority (3A) ...
  ... Passed --   1.0  5   163    3
Test: completion after heal (3A) ...
  ... Passed --   1.1  5    60    3
Test: partitions, one client (3A) ...
  ... Passed --  23.1  5 11383 2047
Test: partitions, many clients (3A) ...
  ... Passed --  22.8  5 16136 2986
Test: restarts, one client (3A) ...
  ... Passed --  19.8  5 11468 2235
Test: restarts, many clients (3A) ...
  ... Passed --  20.0  5 14567 2691
Test: unreliable net, restarts, many clients (3A) ...
  ... Passed --  23.0  5  7155 1000
Test: restarts, partitions, many clients (3A) ...
  ... Passed --  27.4  5 16092 2865
Test: unreliable net, restarts, partitions, many clients (3A) ...
  ... Passed --  28.2  5  5894  558
Test: unreliable net, restarts, partitions, random keys, many clients (3A) ...
  ... Passed --  34.8  7 17566 1387
Test: InstallSnapshot RPC (3B) ...
  ... Passed --   2.1  3   265   63
Test: snapshot size is reasonable (3B) ...
  ... Passed --   1.0  3  2412  800
Test: ops complete fast enough (3B) ...
  ... Passed --   1.2  3  3018    0
Test: restarts, snapshots, one client (3B) ...
  ... Passed --  19.6  5 55028 10952
Test: restarts, snapshots, many clients (3B) ...
  ... Passed --  19.8  5 56273 10255
Test: unreliable net, snapshots, many clients (3B) ...
  ... Passed --  17.9  5  6143  919
Test: unreliable net, restarts, snapshots, many clients (3B) ...
  ... Passed --  20.8  5  6550  892
Test: unreliable net, restarts, partitions, snapshots, many clients (3B) ...
  ... Passed --  27.7  5  6418  658
Test: unreliable net, restarts, partitions, snapshots, random keys, many clients (3B) ...
  ... Passed --  32.5  7 16218 1182
PASS
```
