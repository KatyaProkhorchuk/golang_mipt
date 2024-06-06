package dupcall

import (
    "context"
    "sync"
)

type Call struct {
	mu     sync.Mutex
	result interface{}
	err    error
	done   chan struct{}
	inProgress bool
	currentFunc *func(context.Context) (interface{}, error)
}

func (c *Call) Do(ctx context.Context, cb func(context.Context) (interface{}, error)) (interface{}, error) {
	c.mu.Lock()
	
	if c.inProgress && c.currentFunc != &cb {
		// Если cb уже выполняется, ждем его завершения
		c.mu.Unlock()
		<-c.done
		// Возвращаем результат и ошибку из предыдущего вызова
		c.mu.Lock()
		defer c.mu.Unlock()
		return c.result, c.err
	} else {
        // Устанавливаем флаг в true, что началось выполнение функции
        c.inProgress = true
		if c.currentFunc != &cb {
			c.done = make(chan struct{})
		}
    }
	c.currentFunc = &cb
	c.mu.Unlock()

	// Запускаем cb в отдельной горутине
	go func() {
		defer close(c.done)

		// Вызываем переданную функцию cb
		result, err := cb(ctx)

		// Возвращаем результат и ошибку в канал done
		c.mu.Lock()
		defer c.mu.Unlock()
		if result == nil || c.result != nil {

		}
		// if c.result != result || c.err != err{
		// 	c.done = make(chan struct{})
		// }
		c.result = result
		c.err = err
		c.inProgress = false
	}()

	select {
	case <-ctx.Done():
		// Если контекст отменен, возвращаем ошибку
		return nil, ctx.Err()
	case <-c.done:
		// Если cb завершилась, возвращаем результат и ошибку
		c.mu.Lock()
		defer c.mu.Unlock()
		return c.result, c.err
	}
}