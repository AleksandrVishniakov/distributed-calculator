# Статусы операций
В процессе работы приложения могут использоваться статусы в виде числа
#### Объявление статусов:
```Go
type Status int
const (
	Created Status = iota // 0
	Enqueued // 1
	Calculating // 2
	Finished // 3
	Failed // 4
)
```
#### Значение статусов
* Created (0, *"Создано"*) - статус присваивается при создании нового выражения или задачи
* Enqueued (1, *"В очереди"*) - статус присваивается, задача или выражение находится в очереди на выполнение конкретным агентом
* Calculating (2, *"Выполняется"*) - статус присваивается, когда сервис начинает работу над выражением или задачей
* Finished (3, *"Готово"*) - статус присваивается, работа над выражением или задачей окончена
* Failed (4, *"Ошибка вычисления"*) - статус присваивается, когда в процессе вычисления выражения возникла ошибка (например, деление на ноль)