package main

import (
    "errors"
    "testing"
)

func TestNewIntSet(t *testing.T) {
    slice := []int{10, 20, 10, 30}
    set := NewIntSetFromSlice(slice)
    if len(set)!=3{
        panic(errors.New("NewIntSet gave an unexpected result"))
    }
}

func TestIntSet_IsSubsetOf(t *testing.T) {
    slice1 := []int{10, 20}
    slice2 := []int{10, 20, 10, 30}

    set1 := NewIntSetFromSlice(slice1)
    set2 := NewIntSetFromSlice(slice2)

    if set1.IsSubsetOf(set2)!=true{
        panic(errors.New("IsSubsetOf gave an unexpected result"))
    }
}

func TestIntSet_IsSupersetOf(t *testing.T) {
    slice1 := []int{10, 20, 10, 30}
    slice2 := []int{10, 20}
    set1 := NewIntSetFromSlice(slice1)
    set2 := NewIntSetFromSlice(slice2)

    if set1.IsSupersetOf(set2)!=true{
        panic(errors.New("IsSupersetOf gave an unexpected result"))
    }
}

func TestIntSet_Intersect(t *testing.T) {
    slice1 := []int{10, 20, 10, 30}
    slice2 := []int{10, 20}
    set1 := NewIntSetFromSlice(slice1)
    set2 := NewIntSetFromSlice(slice2)

    intersect := set1.Intersect(set2)

    if len(intersect)!=2{
        panic(errors.New("Intersect gave an unexpected result"))
    }
}