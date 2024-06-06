//go:build !solution

package hogwarts

func findPrereqs(result []string, course string, learnedCourse map[string]bool, prereqsCourse map[string]bool, prereqs map[string][]string) []string {
	// рекурсивно просматриваем курсы
	if learnedCourse[course] { // если курс изучен
		return result
	}
	if prereqsCourse[course] { // возник цикл
		panic("error")
	}
	prereqsCourse[course] = true
	for _, prc := range prereqs[course] {
		result = findPrereqs(result, prc, learnedCourse, prereqsCourse, prereqs)
	}
	prereqsCourse[course] = false
	learnedCourse[course] = true
	result = append(result, course)
	return result
}

func GetCourseList(prereqs map[string][]string) []string {
	result := make([]string, 0)
	learnedCourse := make(map[string]bool)
	prereqsCourse := make(map[string]bool)
	for course := range prereqs {
		result = findPrereqs(result, course, learnedCourse, prereqsCourse, prereqs)
	}
	return result
}
