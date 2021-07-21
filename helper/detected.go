package helper

import (
	"sort"

	"github.com/Kagami/go-face")


type Detected struct {
	Name     string
	distance float64
}

type DetectedSlice []Detected

func (a DetectedSlice) Len() int           { return len(a) }
func (a DetectedSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DetectedSlice) Less(i, j int) bool { return a[i].distance < a[j].distance }

func (slice *DetectedSlice) FillSortDetected(uFaceDescriptor face.Descriptor, samples []face.Descriptor, labels []string) {
	for k, v := range samples {
		dist := face.SquaredEuclideanDistance(v, uFaceDescriptor)
		if (*slice).Len() > 0 && (*slice)[len(*slice)-1].Name == labels[k] {
			continue
		}
		if dist > 0.6 {
			break
		}

		*slice = append(*slice, Detected{Name: labels[k], distance: dist})
	}
	*slice = append(*slice, Detected{Name: "Unknown", distance: 0.6})
	sort.Sort(*slice)

}