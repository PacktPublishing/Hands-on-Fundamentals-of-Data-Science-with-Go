func (set StringSet) IsEqualSubset(set2 StringSet) bool {
    if set.IsSubsetOf(set2) && set2.IsSubsetOf(set) {
        return true
    }
    return false
}

func (set StringSet) IsProperSubsetOf(set2 StringSet) bool {
    for key1 := range set {
        if !set2.Contains(key1) {
            return false
        }
    }
    if set.IsEqualSubset(set2) {
        return false
    }
    return true
}