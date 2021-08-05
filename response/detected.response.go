package response

import (
	"goface-api/models"
	"sort"

	"github.com/Kagami/go-face"
)

type Detected struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Distance float64 `json:"distance"`
}

type DetectedSlice []Detected

func (a DetectedSlice) Len() int           { return len(a) }
func (a DetectedSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DetectedSlice) Less(i, j int) bool { return a[i].Distance < a[j].Distance }

func (slice *DetectedSlice) FillSortDetectedFromDB(udesc face.Descriptor, sampleFaces []models.Face, threshold float64) {
	*slice = append(*slice, Detected{Name: "Unknown", Distance: threshold})
	for _, value := range sampleFaces {
		for _, desc := range value.Descriptors {
			dist := face.SquaredEuclideanDistance(udesc, desc)
			if dist > threshold {
				continue
			}

			if lastSlice := (*slice)[len(*slice)-1]; lastSlice.Name == value.Name {
				if dist < lastSlice.Distance {
					lastSlice.Distance = dist
				}
				continue
			}

			*slice = append(*slice, Detected{Id: value.Id, Name: value.Name, Distance: dist})
		}

	}
	sort.Sort(slice)

}

func (slice *DetectedSlice) FillSortDetected(udescriptor face.Descriptor, samples []face.Descriptor, labels []string, threshold float64) {
	*slice = append(*slice, Detected{Name: "Unknown", Distance: threshold})
	for k, v := range samples {
		dist := face.SquaredEuclideanDistance(udescriptor, v)
		if dist > threshold {
			continue
		}
		lastSlice := &(*slice)[len(*slice)-1]
		if lastSlice.Name == labels[k] {
			if dist < lastSlice.Distance {
				lastSlice.Distance = dist
			}
			continue
		}

		*slice = append(*slice, Detected{Name: labels[k], Distance: dist})
	}
	sort.Sort(*slice)

}
