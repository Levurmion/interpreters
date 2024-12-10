package lr1closureset

import (
	"interpreters/internal/parser/lr1item"
	"interpreters/utilities/arrays"
	"interpreters/utilities/maps"
	"interpreters/utilities/sets"
)

type LR1ClosureSet struct {
	kernelItems map[string]*lr1item.LR1Item
	items map[string]*lr1item.LR1Item
}

func NewEmptyLR1ClosureSet() *LR1ClosureSet {
	kernelItems := make(map[string]*lr1item.LR1Item)
	items := make(map[string]*lr1item.LR1Item)
	return &LR1ClosureSet{kernelItems, items}
}

// Creates a new `LR1ClosureSet` where all initially loaded items are considered
// `kernel` items.
func NewLR1ClosureSet(LR1Items ...*lr1item.LR1Item) *LR1ClosureSet {
	kernelItems := make(map[string]*lr1item.LR1Item)
	items := make(map[string]*lr1item.LR1Item)
	for _, item := range LR1Items {
		itemName := item.GetName()
		if _, exists := kernelItems[itemName]; !exists {
			kernelItems[itemName] = item
			items[itemName] = item
		}
	}
	return &LR1ClosureSet{kernelItems, items}
}

func (cs *LR1ClosureSet) Add(LR1Item *lr1item.LR1Item) {
	itemName := LR1Item.GetName()
	if _, exists := cs.items[itemName]; !exists {
		cs.items[itemName] = LR1Item
	}
}

func (cs *LR1ClosureSet) Delete(LR1Item *lr1item.LR1Item) {
	itemName := LR1Item.GetName()
	delete(cs.items, itemName)
}

func (cs *LR1ClosureSet) GetItems() []*lr1item.LR1Item {
	items := []*lr1item.LR1Item{}
	for _, item := range cs.items {
		items = append(items, item)
	}
	return items
}

func (cs *LR1ClosureSet) Clone() *LR1ClosureSet {
	kernelCopy := maps.DeepCopy(cs.kernelItems)
	itemsCopy := maps.DeepCopy(cs.items)
	return &LR1ClosureSet{kernelCopy, itemsCopy}
}

// Compares two `LR1ClosureSets` by both kernel and lookahead identity.
func (thisSet *LR1ClosureSet) IsEqual(otherSet LR1ClosureSet) bool {
	thisKernel := thisSet.kernelItems
	otherKernel := otherSet.kernelItems

	// check that both kernels are the same size
	if len(thisKernel) != len(otherKernel) {
		return false
	}

	// check that all items exist and lookahead sets are equal
	for itemName, thisItem := range thisKernel {
		otherItem, exists := otherKernel[itemName]
		if !exists {
			return false
		}
		if !thisItem.LookaheadSet.IsEqual(otherItem.LookaheadSet) {
			return false
		}
	}

	return true
}

// Retrieve items whose productions are complete.
func (cs *LR1ClosureSet) GetCompletedItems() []*lr1item.LR1Item {
	items := cs.GetItems()
	return arrays.Filter(items, func (item *lr1item.LR1Item) bool {
		return item.ProductionIsComplete()
	})
}

// Get all symbols that can follow the `dot` of every item in
// the `LR1ClosureSet`.
func (cs *LR1ClosureSet) GetTransitionSymbols() sets.Set[string] {
	items := cs.GetItems()
	transitionSymbols := arrays.Map(items, func (item *lr1item.LR1Item) string {
		return item.GetNextSymbol()
	})
	return sets.NewSet(transitionSymbols...)
}