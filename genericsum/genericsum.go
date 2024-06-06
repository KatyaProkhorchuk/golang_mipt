//go:build !solution

package genericsum

import (
	"math/cmplx"
	"sync"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// минимум из 2 переменных
func Min[T constraints.Ordered](a, b T) T {
	if a > b {
		return b
	}
	return a
}

// сортировка слайса inplace
func SortSlice[T constraints.Ordered](a []T) {
	slices.Sort(a)
}

// равенство 2 мап. Значения мап сравниваются через обычный оператор =
func MapsEqual[Map1, Map2 ~map[KEY]VAL, KEY, VAL comparable](a Map1, b Map2) bool {
	if len(a) != len(b) {
		return false
	}
	for k, val_a := range a {
		if val_b, ok := b[k]; !ok || val_b != val_a {
			return false
		}
	}
	return true
}

// содержит ли слайс заданный элемент
func SliceContains[T comparable](s []T, v T) bool {
	for _, val := range s {
		if val == v {
			return true
		}
	}
	return false
}

// сделать из нескольких каналов один
func MergeChans[T any](chs ...<-chan T) <-chan T {
	result := make(chan T)
	go func() {
		defer close(result) // Закрыть канал result после завершения работы горутины
		var wg sync.WaitGroup
		wg.Add(len(chs))

		for _, c := range chs {
			go func(c <-chan T) {
				defer wg.Done() // Уменьшить счетчик горутин при завершении
				for v := range c {
					result <- v
				}
			}(c)
		}

		wg.Wait() // Ждать завершения всех горутин, обрабатывающих каналы chs
	}()
	return result
}

type Matrix interface {
	constraints.Integer | constraints.Complex | constraints.Float
}

// проверка, является ли квадратная матрица Эрмитовой.
func IsHermitianMatrix[T Matrix](m [][]T) bool {
	h := len(m)
	w := len(m[0])
	if h < 1 || w < 1 {
		return true
	}
	if h != w {
		return false
	}
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			switch any(m[i][j]).(type) {
			case complex64:
				transp, ok := any(m[j][i]).(complex64)
				if !ok {
					return false
				}
				curr, ok := any(m[i][j]).(complex64)
				if !ok {
					return false
				}
				if real(transp) != real(curr) || imag(transp) != -imag(curr) {
					return false
				}
			case complex128:
				transp, ok := any(m[j][i]).(complex128)
				if !ok {
					return false
				}
				curr, ok := any(m[i][j]).(complex128)
				if !ok {
					return false
				}
				if cmplx.Conj(curr) != transp {
					return false
				}
			default:
				if m[i][j] != m[j][i] {
					return false
				}
			}
		}
	}
	return true
}
