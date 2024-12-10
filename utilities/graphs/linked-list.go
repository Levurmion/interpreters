package graphs

type ListNode[T any] struct {
	Item T
	Next *ListNode[T]
	Prev *ListNode[T]
}

/*
The `LinkedList` can be used in place of built-in arrays for data structures that require O(1)
access to items at the start and end of the list.
*/
type LinkedList[T any] struct {
	length int
	dummyHead *ListNode[T]
	dummyTail *ListNode[T]
}

func NewLinkedList[T any]() *LinkedList[T] {
	dummyHead := new(ListNode[T])
	dummyTail := new(ListNode[T])
	dummyHead.Next = dummyTail
	dummyTail.Prev = dummyHead
	return &LinkedList[T]{
		0,
		dummyHead,
		dummyTail,
	}
}

func (list *LinkedList[T]) GetHead() *ListNode[T] {
	head := list.dummyHead.Next
	return head
}

func (list *LinkedList[T]) GetTail() *ListNode[T] {
	tail := list.dummyTail.Prev
	return tail
}

func (list *LinkedList[T]) GetLength() int {
	return list.length
}

// Add item to the tail of the `LinekdList`
func (list *LinkedList[T]) Append(item T) {
	newNode := new(ListNode[T])
	newNode.Item = item
	currTail := list.GetTail()

	// stitch newNode between the current tail and dummyTail
	currTail.Next = newNode
	newNode.Prev = currTail
	list.dummyTail.Prev = newNode
	newNode.Next = list.dummyTail
	list.length++
}

// Add item to the head of the `LinkedList`
func (list *LinkedList[T]) Unshift(item T) {
	newNode := new(ListNode[T])
	newNode.Item = item
	currHead := list.GetHead()

	// stitch newNode between the current head and dummyHead
	currHead.Prev = newNode
	newNode.Next = currHead
	list.dummyHead.Next = newNode
	newNode.Prev = list.dummyHead
	list.length++
}

// Remove and return the last node of the `LinkedList`
func (list *LinkedList[T]) Pop() *ListNode[T] {
	currTail := list.GetTail()
	if currTail == list.dummyHead {
		return nil
	}

	nextTail := currTail.Prev
	dummyTail := currTail.Next

	// remove references to list in currTail
	currTail.Next = nil
	currTail.Prev = nil

	// stitch nextTail into dummyTail
	nextTail.Next = dummyTail
	dummyTail.Prev = nextTail
	list.length--

	return currTail
}

// Remove and return the first node of the `LinkedList`
func (list *LinkedList[T]) Shift() *ListNode[T] {
	currHead := list.GetHead()
	if currHead == list.dummyTail {
		return nil
	}

	nextHead := currHead.Next
	dummyHead := currHead.Prev

	// remove references to list in currHead
	currHead.Next = nil
	currHead.Prev = nil

	// stitch nextHead into dummyHead
	nextHead.Prev = dummyHead
	dummyHead.Next = nextHead
	list.length--

	return currHead
}

