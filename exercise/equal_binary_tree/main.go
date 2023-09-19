package main

import (
	"fmt"
	"tour_of_go/exercise/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
    var dfs func(t *tree.Tree, ch chan int)
    dfs = func(t *tree.Tree, ch chan int) {
        if t == nil {
            return
        }
        dfs(t.Left, ch)
        ch <- t.Value
        dfs(t.Right, ch)
    }
    dfs(t, ch)
    close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
    ch1 := make(chan int)
    ch2 := make(chan int)
    go Walk(t1, ch1)
    go Walk(t2, ch2)

    arr1 := make([]int, 0)
    for num1 := range ch1 {
        arr1 = append(arr1, num1)
    }

    arr2 := make([]int, 0)
    for num2 := range ch2 {
        arr2 = append(arr2, num2)
    }

    if len(arr1) != len(arr2) {
        return false
    }

    for k, v := range arr1 {
        if v != arr2[k] {
            return false
        }
    }
    return true
}

func main() {
    ch := make(chan int)
    go Walk(tree.New(10), ch)
    for num := range ch {
        fmt.Println(num)
    }

    isSame := Same(tree.New(1), tree.New(2))
    fmt.Println(isSame)

    isSame = Same(tree.New(1), tree.New(1))
    fmt.Println(isSame)
}
