# Select Priority Assumption

Этот пример показывает частый channel scheduling bug:

1. У dispatcher есть high-priority и low-priority queues.
2. Обе очереди могут быть ready одновременно.
3. Код считает, что `select` выбирает первый ready case.

Go `select` не дает такой priority guarantee. Если несколько cases ready, выбирается один pseudo-random case.

Нарушение находится в [internal/dispatcher/dispatcher.go](internal/dispatcher/dispatcher.go): `Next` ставит high-priority receive первым и считает, что этого достаточно.

## Запуск

```sh
make run
```

Пример вывода:

```text
select does not prioritize the first ready case
run make lint to see the repeated schedule check
one run selected low job batch-1
```

Последняя строка может упоминать и `high`, и `low`. В этом и смысл: один запуск не доказывает schedule assumption.

## Обычный Тест

```sh
make test
```

Обычный тест делает ready только high-priority queue. Он проходит, но не покрывает рискованное состояние, где обе очереди ready.

## Как Поймать Через Repeated Schedule Check

```sh
make lint
```

Target `lint` запускает bug-revealing test. Он много раз создает оба ready cases и падает, как только `select` выбирает low-priority queue.

Ожидаемый отчет:

```text
--- FAIL: TestSelectDoesNotGuaranteePriority (0.00s)
    main_test.go:29: select chose low-priority job while high-priority job was ready
FAIL
FAIL	github.com/dvordrova/find_bugs/concurrency/select_priority_assumption Xs
?   	github.com/dvordrova/find_bugs/concurrency/select_priority_assumption/internal/dispatcher	[no test files]
FAIL
```

Это легкая локальная версия testing idea из GFuzz: менять порядок, в котором concurrent messages становятся observable, и проверять, остается ли программа корректной.

## Одно Исправление

Сделать явную non-blocking проверку high-priority queue перед low-priority queue:

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

Второй `select` все еще позволяет dispatcher ждать, когда work не ready, но больше не считает low-priority work равным, когда high-priority work уже доступен.
