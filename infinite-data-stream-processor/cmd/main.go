package main

import (
	"errors"
	"fmt"
	"iter"
)

type Page struct {
	Items []string
	Next  int
}

// fetchPage simulates a paginated API. Returns io.EOF-style sentinel when done.
func fetchPage(cursor int) (Page, error) {
	pages := []Page{
		{Items: []string{"item-1", "item-2", "item-3"}, Next: 1},
		{Items: []string{"item-4", "item-5"}, Next: 2},
		{Items: []string{"item-6", "item-7", "item-8", "item-9"}, Next: 3},
	}

	if cursor >= len(pages) {
		return Page{}, errors.New("no more pages")
	}
	return pages[cursor], nil
}

// Pages returns a push iterator (iter.Seq2) that yields each item from all
// pages alongside any fetch error. The caller breaks out of the range loop on
// error; the iterator stops yielding the moment yield returns false.
func Pages(startCursor int) iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		cursor := startCursor
		index := 0
		for {
			page, err := fetchPage(cursor)
			if err != nil {
				return
			}
			for _, item := range page.Items {
				if !yield(index, item) {
					return
				}
				index++
			}
			cursor = page.Next
		}
	}
}

func main() {
	fmt.Println("--- for...range (push iterator) ---")
	for i, item := range Pages(0) {
		fmt.Printf("[%d] %s\n", i, item)
	}

	fmt.Println("\n--- iter.Pull2 (manual stepping) ---")
	next, stop := iter.Pull2(Pages(0))
	defer stop()

	// Step through manually — pull only the first 5 items then stop early.
	for {
		i, item, ok := next()
		if !ok {
			break
		}
		fmt.Printf("[%d] %s\n", i, item)
		//if i == 4 {
		//	fmt.Println("(stopped early via Pull2)")
		//	break
		//}
	}
}
