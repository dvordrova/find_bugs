# Select Priority Assumption

This example shows a common channel scheduling bug:

1. A dispatcher has high-priority and low-priority queues.
2. Both queues can be ready at the same time.
3. The code assumes that a `select` chooses the first ready case.

Go `select` does not provide that priority guarantee. If multiple cases are ready, one is selected pseudo-randomly.

The violation is in [internal/dispatcher/dispatcher.go](internal/dispatcher/dispatcher.go): `Next` places the high-priority receive first and assumes that is enough.

## Run

```sh
make run
```

Example output:

```text
select does not prioritize the first ready case
run make lint to see the repeated schedule check
one run selected low job batch-1
```

The last line can mention either `high` or `low`. That is the point: a single run does not prove the schedule assumption.

## Ordinary Test

```sh
make test
```

The ordinary test only makes the high-priority queue ready. It passes, but it does not cover the risky state where both queues are ready.

## Catch With Repeated Schedule Check

```sh
make lint
```

The `lint` target runs a bug-revealing test. It repeatedly creates both ready cases and fails as soon as `select` chooses the low-priority queue.

Expected report:

```text
--- FAIL: TestSelectDoesNotGuaranteePriority (0.00s)
    main_test.go:29: select chose low-priority job while high-priority job was ready
FAIL
FAIL	github.com/dvordrova/find_bugs/concurrency/select_priority_assumption Xs
?   	github.com/dvordrova/find_bugs/concurrency/select_priority_assumption/internal/dispatcher	[no test files]
FAIL
```

This is a lightweight, local version of the testing idea behind GFuzz: change the order in which concurrent messages become observable and check whether the program still behaves correctly.

## One Fix

Use an explicit non-blocking high-priority check before considering the low-priority queue:

```go
select {
case job := <-highPriority:
	return job
default:
}

select {
case job := <-highPriority:
	return job
case job := <-lowPriority:
	return job
}
```

The second `select` still lets the dispatcher block when no work is ready, but it no longer treats low-priority work as equal when high-priority work is already available.
