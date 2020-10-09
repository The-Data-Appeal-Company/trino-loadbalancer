package discovery

func containsAllTags(src, test map[string]string) bool {
	for k, v := range test {
		sv, present := src[k]
		if !present || v != sv {
			return false
		}
	}
	return true
}
