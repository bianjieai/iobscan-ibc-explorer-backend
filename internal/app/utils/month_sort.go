package utils

type MonthSlice []string

func (m MonthSlice) Len() int {
	return len(m)
}

func (m MonthSlice) Less(i, j int) bool {
	dateI := StringToMonthDate(m[i])
	dateJ := StringToMonthDate(m[j])
	return dateI.Before(dateJ)
}

func (m MonthSlice) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
