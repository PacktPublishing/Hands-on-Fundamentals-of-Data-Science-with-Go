package main

type IntSet map[int]bool

func NewIntSet()IntSet{
    set := IntSet{}
    return set
}

func NewIntSetFromSlice(slice []int) IntSet {
    set := IntSet{}
    for _, val := range slice{
        set[val] = true
    }
    return set
}

func (set IntSet)Add(val int) IntSet {
    set[val]=true
    return set
}

func (set IntSet)Contains(val int) bool {
    if _, ok := set[val]; ok{
        return true
    }
    return false
}

func (set IntSet)IsSubsetOf(set2 IntSet) bool {
    for key1 := range set {
        if !set2.Contains(key1){
            return false
        }
    }
    return true
}

func (set IntSet)IsSupersetOf(set2 IntSet) bool {
    for key2 := range set2{
        if !set.Contains(key2){
            return false
        }
    }
    return true
}


func (set IntSet)Intersect(set2 IntSet) IntSet {
    newSet := NewIntSet()

    for key1 := range set {
        for key2 := range set2{
            if key1 == key2{
                newSet.Add(key1)
            }
        }
    }
    return newSet
}