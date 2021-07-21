package helper

type Detected struct{
	Name string
	distance float64
}

type DetectedSlice []Detected

func (a DetectedSlice) Len() int           { return len(a) }
func (a DetectedSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DetectedSlice) Less(i, j int) bool { return a[i].distance < a[j].distance }