package utils

func Contains(s []int, val int) bool {
	for _, item := range s {
		if item == val {
			return true
		}
	}
	return false
}

func Remove(slice []int, s int) []int {
	for i, num := range slice {
		if num == s {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return slice
}

func CommonId(user1 []int, user2 []int) []int {
	seen := make([]int, 1000)
	for i := range user1 {
		seen[user1[i]]++
	}

	res := make([]int, 0)
	for i := range user2 {
		if seen[user2[i]] > 0 {
			res = append(res, user2[i])
			seen[user2[i]] = 0
		}
	}
	return res
}
