package redirects

import "fmt"

type Token struct {
	token string
	index int
	next  *Token
	prev  *Token
}

type List struct {
	head *Token
	tail *Token
}

func (L *List) Insert(token string, index int) {
	list := &Token{
		next:  L.head,
		token: token,
		index: index,
	}
	if L.head != nil {
		L.head.prev = list
	}
	L.head = list

	l := L.head
	for l.next != nil {
		l = l.next
	}
	L.tail = l
}

func (l *List) Print() {
	list := l.head
	for list != nil {
		fmt.Printf("%d%+v ->", list.index, list.token)
		list = list.next
	}
	fmt.Println()
}

func PrintListFrom(list *Token) {
	for list != nil {
		fmt.Printf("%d%v ->", list.index, list.token)
		list = list.next
	}
	fmt.Println()
}

func (l *List) Reverse() {
	curr := l.head
	var prev *Token
	l.tail = l.head

	for curr != nil {
		next := curr.next
		curr.next = prev
		prev = curr
		curr = next
	}
	l.head = prev
	// PrintListFrom(l.head)
}
