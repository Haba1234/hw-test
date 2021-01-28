package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	// Place your code here
	Len() int                          // длина списка
	Front() *ListItem                  // первый элемент списка
	Back() *ListItem                   // последний элемент списка
	PushFront(v interface{}) *ListItem // добавить значение в начало
	PushBack(v interface{}) *ListItem  // добавить значение в конец
	Remove(i *ListItem)                // удалить элемент
	MoveToFront(i *ListItem)           // переместить элемент в начало
}

type ListItem struct {
	// Place your code here
	Value interface{} // значение
	Next  *ListItem   // следующий елемент
	Prev  *ListItem   // предыдущий элемент
}

type list struct {
	elements []*ListItem
}

func NewList() List {
	return new(list)
}

// вычисление длины списка
func (l *list) Len() int {
	return len(l.elements)
}

// вернуть адрес первого элемента списка или nil, если список пуст
func (l *list) Front() *ListItem {
	if result := l.Len(); result > 0 {
		return l.elements[0]
	}
	return nil
}

// вернуть адрес последнего элемента списка или nil, если список пуст
func (l *list) Back() *ListItem {
	if result := l.Len(); result > 0 {
		return l.elements[result-1]
	}
	return nil
}

// добавить значение в начало списка
func (l *list) PushFront(v interface{}) *ListItem {
	// первая запись? создать певвый элемент списка
	if l.Front() == nil {
		l.elements = append(l.elements, &ListItem{Value: v, Next: nil, Prev: nil})
		return l.Front()
	}

	// добавим новый элемент в начало списка
	newSlice := []*ListItem{{Value: v, Next: l.Front(), Prev: nil}} // создали слайс с новым элементом списка
	l.elements = append(newSlice, l.elements...)                    // пееенесли в новый слайс все элементы старого
	l.elements[1].Prev = l.Front()                                  // обновили адрес в Prev для следующего элемента списка
	return l.Front()
}

// добавить значение в конец списка
func (l *list) PushBack(v interface{}) *ListItem {
	// проеерка, что список не пустой
	if prev := l.Back(); prev != nil {
		// добавить значение в конец
		l.elements = append(l.elements, &ListItem{Value: v, Next: nil, Prev: prev})
		l.elements[l.Len()-2].Next = l.Back() // у предыдущей записи обновить ссылку Next на новый элемент списка
		return l.Back()
	}

	// список пуст, создать первую запись
	return l.PushFront(v)
}

// удалить элемент списка
func (l *list) Remove(i *ListItem) {
	// список пуст!
	listLen := l.Len()
	if listLen == 0 {
		return
	}
	// проверка, элемент последний в списке?
	if l.Back().Value == i.Value {
		l.elements = l.elements[:listLen-1] // обрезаем слайс на крайний элемент
		if l.Len() > 0 {
			// обнуляем ссылку у теперь уже нового крайнего элемента, если он еще остался
			l.Back().Next = nil
		}
		return
	}
	// удаляемый элемент не последний в списке. Ищем где он...
	for num, val := range l.elements {
		if val.Value == i.Value {
			l.elements = append(l.elements[:num], l.elements[num+1:]...) // удаление
			// обновление ссылкок у предыдущего и следующего элементов списка
			if num > 0 {
				l.elements[num-1].Next = l.elements[num]
				l.elements[num].Prev = l.elements[num-1]
			} else {
				// если удаляемый элемент первый, только обнулить ссылку у следующего
				l.elements[num].Prev = nil
			}
			return
		}
	}
}

// переместить элемент в начало списка
func (l *list) MoveToFront(i *ListItem) {
	// ищем элемент для перемещения в списке...
	for num, val := range l.elements {
		// проверяем, что он не первый
		if (val.Value == i.Value) && (num > 0) {
			l.Remove(i)          // удаляем из текущей позиции
			l.PushFront(i.Value) // и вставляем значение в начало списка
			return
		}
	}
}
