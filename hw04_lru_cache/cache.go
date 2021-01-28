package hw04_lru_cache //nolint:golint,stylecheck
import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу.
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу.
	Clear()                              // Очистить кэш.
}

type lruCache struct {
	mx       sync.RWMutex
	capacity int               // ёмкость (количество сохраняемых в кэше элементов).
	queue    List              // очередь [последних используемых элементов] на основе двусвязного списка.
	items    map[Key]*ListItem // словарь, отображающий ключ (строка) на элемент очереди.
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Добавить значение в кэш по ключу.
func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.mx.Lock()         // включаем защиту при параллельном доступе.
	defer lru.mx.Unlock() // выключаем защиту перед выходом из функции.

	objValue := cacheItem{key: key, value: value}
	_, ok := lru.items[key]
	lru.items[key] = &ListItem{Value: objValue} // обновляем значение в словаре или добавляем в словарь.
	if !ok {
		// если элемента нет в словаре, то добавить в словарь и в начало очереди
		// (при этом, если размер очереди больше ёмкости кэша, то необходимо
		// удалить последний элемент из очереди и его значение из словаря).
		if lru.queue.Len() == lru.capacity { // очередь достигла максимума.
			valBack := lru.queue.Back() // крайний элемент в списке.

			lru.queue.Remove(valBack)                        // удалить крайний элемент из списка.
			delete(lru.items, valBack.Value.(cacheItem).key) // удалить элемент из словаря.
		}
		lru.queue.PushFront(objValue) // добавить элемент в начало списка.
		return false
	}

	// если элемент присутствует в словаре, то обновить его значение
	// и переместить элемент в начало очереди.
	lru.queue.MoveToFront(lru.items[key])
	return true
}

// Получить значение из кэша по ключу.
func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.mx.RLock()         // включаем защиту при параллельном доступе.
	defer lru.mx.RUnlock() // выключаем защиту перед выходом из функции.

	if elem, ok := lru.items[key]; ok {
		// если элемент присутствует в словаре, то переместить элемент
		// в начало очереди и вернуть его значение и true.
		lru.queue.MoveToFront(lru.items[key])
		return elem.Value.(cacheItem).value, true
	}
	// если элемента нет в словаре, то вернуть nil и false.
	return nil, false
}

// Очистить кэш.
func (lru *lruCache) Clear() {
	lru.mx.Lock()         // включаем защиту при параллельном доступе.
	defer lru.mx.Unlock() // выключаем защиту перед выходом из функции.

	lru.items = make(map[Key]*ListItem, lru.capacity) // очистка кэша.
	for i := 0; i < lru.queue.Len(); i++ {            // очистка очереди.
		lru.queue.Remove(lru.queue.Back())
	}
}
