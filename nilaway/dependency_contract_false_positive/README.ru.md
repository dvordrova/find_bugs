# SDK Init Global False Positive

Этот пример показывает false positive NilAway, который уходит в dependency module.

Приложение зависит от маленького локального SDK-модуля:

```go
require github.com/acme/contractsdk v0.0.0

replace github.com/acme/contractsdk => ./contractsdk
```

SDK экспортирует package-level pointer:

```go
var DefaultPlan *Plan

func init() {
	DefaultPlan = &Plan{
		Name: "enterprise",
	}
}
```

В runtime это безопасно: Go запускает package `init` до `main`, поэтому `contractsdk.DefaultPlan` инициализирован до того, как приложение его читает.

## Запуск

```sh
make run
```

Ожидаемый результат:

```text
tenant uses enterprise plan
```

## Как поймать через NilAway

```sh
make lint
```

Ожидаемый вывод:

```text
main.go:16:38: Potential nil panic detected. Observed nil flow from source to dereference point:
	- contractsdk/tenant.go:8:5: nilable value assigned into global variable `DefaultPlan`
	- dependency_contract_false_positive/main.go:12:12: global variable `DefaultPlan` passed as arg `plan` to `printPlan()` via the assignment(s):
		- `contractsdk.DefaultPlan` to `plan` at dependency_contract_false_positive/main.go:10:2
	- dependency_contract_false_positive/main.go:16:38: function parameter `plan` accessed field `Name` (nilaway)
	fmt.Printf("tenant uses %s plan\n", plan.Name)
	                                    ^
```

Читать отчет надо как flow от возможного nil source к dereference:

1. `contractsdk/tenant.go:8:5`: `DefaultPlan` - package-level pointer. Его zero value равен nil.
2. `main.go:10:2`: приложение присваивает `contractsdk.DefaultPlan` локальной переменной `plan`.
3. `main.go:12:12`: приложение передает `plan` в `printPlan`.
4. `main.go:16:38`: `printPlan` разыменовывает `plan.Name`.

False-positive часть в том, что NilAway не доказывает: SDK `init` всегда присваивает `DefaultPlan` до того, как `main` его использует. Runtime порядок initialization гарантирует, но static analysis осторожничает вокруг mutable global pointers.

## Обработать подтвержденный false positive

В production не стоит прятать это широким dependency exclusion. Global pointer сам по себе рискованная форма API: другая версия SDK, тест или будущая мутация могут сделать его nil.

Оба lint config держат `include-pkgs` явным:

```yaml
include-pkgs: github.com/dvordrova/find_bugs/nilaway/dependency_contract_false_positive,github.com/acme/contractsdk
```

Это значит, что NilAway анализирует приложение и SDK contract, который участвует в этом примере. Не надо класть туда все SDK по умолчанию. Добавляй packages, которыми ты владеешь, или packages, nil contracts которых ты осознанно хочешь проверять через NilAway. Пустой `include-pkgs` удобен для исследования, но для CI обычно слишком широкий: third-party и generated packages могут добавить шум и замедлить прогон.

В этом примере используется узкий golangci-lint exclusion в [.golangci.fixed.yaml](.golangci.fixed.yaml):

```yaml
linters:
  exclusions:
    warn-unused: true
    rules:
      # Known false positive: contractsdk.DefaultPlan is an SDK global pointer
      # initialized from contractsdk.init before main runs. Keep this narrow:
      # only NilAway, only this app file, only the DefaultPlan global-flow report.
      # In production, back this kind of suppression with an SDK contract or test.
      - linters:
          - nilaway
        path: ^main\.go$
        text: nilable value assigned into global variable `DefaultPlan`
```

Он специально точечный:

1. `linters: [nilaway]` оставляет все остальные линтеры включенными.
2. `path: ^main\.go$` ограничивает правило этим app entrypoint.
3. `text: ...DefaultPlan` матчится на известный SDK global-pointer false positive, а не на любой nil dereference.
4. `warn-unused: true` заставит golangci-lint предупредить, если exclusion перестал матчиться после изменения кода.

Запуск:

```sh
make lint-fixed
```

или напрямую:

```sh
./custom-gcl run --config .golangci.fixed.yaml
```

Ожидаемый результат:

```text
0 issues.
```

Custom linter config фиксирует версию NilAway module в [.custom-gcl.yml](.custom-gcl.yml), поэтому `make lint` и `make lint-fixed` используют одну и ту же версию analyzer, пока ее явно не обновят.

`make tool-update` - maintainer-команда для осознанного обновления tool dependency в `go.mod`; для запуска примера она не нужна.

## Production alternatives

В реальном коде часто лучше вообще не протаскивать nil-capable SDK global глубоко в business logic. Проверить SDK boundary один раз:

```go
func defaultPlan() *contractsdk.Plan {
	plan := contractsdk.DefaultPlan
	if plan == nil {
		log.Fatal("SDK default plan is not initialized")
	}

	return plan
}
```

Если SDK под твоим контролем, лучше вообще не экспортировать mutable pointer global:

```go
func DefaultPlan() *Plan
```

или:

```go
func DefaultPlan() (*Plan, error)
```
