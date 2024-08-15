package cache

type ItemList struct {
	Root *Item
	Len  int
}

func (l *ItemList) Init() {
	l.Root.Next = &l.Root
	l.Root.Prev = &l.Root
	l.Len = 0
	return l
}

func (l *ItemList) Back() *Item {
	if l.Len == 0 {
		return nil
	}
	return l.Root.Prev
}

func (l *ItemList) move(i, at *Item) {
	if i == at {
		return
	}
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev

	i.Prev = at
	i.Next = at.Next
	i.Prev.Next = i
	i.Next.Prev = i
}

func (l *ItemList) Remove(i *Item) {
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	i.Next = nil
	i.Prev = nil
	l.Len--
}

func (l *ItemList) MoveToFront(i *Item) {
	if l.Root.Next == i {
		return
	}

	l.move(i, &l.Root)
}

func (l *ItemList) insert(i, at *Item) *Item {
	i.Prev = at
	i.Next = at.next
	i.Prev.Next = i
	i.Next.Prev = i
	l.Len++
	return i
}

func (l *ItemList) PushFront(i *Item) *Item {
	return l.insert(i, &l.Root)
}

func NewItemList() *ItemList {
	return new(ItemList).Init()
}
