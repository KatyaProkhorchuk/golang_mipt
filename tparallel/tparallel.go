//go:build !solution

package tparallel

type T struct {
	isParallel bool // параллельный ли тест
	parent *T // ссылка на родительский тест
	subtest []*T // для подтестов
	done chan struct{} // что бы сообщить что тест завершен
	barrier chan struct{} // для синхронизации горутин
}


func (t *T) Parallel() {
	if t.isParallel {
		panic("double parallel")
	}
	t.isParallel = true
	t.parent.subtest = append(t.parent.subtest, t)
	t.done <- struct{}{}
	<-t.parent.barrier

}
func (t *T) Run(subtest func(t *T)) {
	// для вызыываемого создадим объект
	subtestRun := &T{
		parent : t,
		subtest: make([]*T, 0),
		done: make(chan struct{}),
		barrier: make(chan struct{}),
	}
	go func() {
		subtest(subtestRun) // запускаем тест
		if len(subtestRun.subtest) > 0 { // если есть вложенные тесты
			close(subtestRun.barrier) // что бы подтесты смогли выполниться
			for _, test := range subtestRun.subtest {
				<-test.done // ждем выполнения всех вложенных тестов
			}
		}
		if subtestRun.isParallel {
			subtestRun.parent.done <- struct{}{} // сообщаем радителю о завершении
		}
		subtestRun.done <- struct{}{} 
	}()
	<-subtestRun.done
}

func Run(topTests []func(t *T)) {
	topTestsRun := &T{
		parent : nil,
		subtest: make([]*T, 0),
		done: make(chan struct{}),
		barrier: make(chan struct{}),
	}
	for _, test := range topTests {
		topTestsRun.Run(test)
	}
	close(topTestsRun.barrier)
	if len(topTestsRun.subtest) > 0 {
		<-topTestsRun.done
	}
}
