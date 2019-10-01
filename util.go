package cloudvision

import (
	"log"
	"reflect"
	"sort"
	"strings"

	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

func getField(v *pb.Vertex, field string) int {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return int(f.Int())
}

func average(in []int) float64 {
	var sum int
	for _, item := range in {
		sum += item
	}

	var avg float64
	avg = float64(sum) / float64(len(in))
	return avg
}

func findByVertices(eas []*pb.EntityAnnotation, vertices []*pb.Vertex) (*pb.EntityAnnotation, int) {
	for i, item := range eas {
		if item.BoundingPoly.Vertices[0] == vertices[0] && item.BoundingPoly.Vertices[1] == vertices[1] &&
			item.BoundingPoly.Vertices[2] == vertices[2] && item.BoundingPoly.Vertices[3] == vertices[3] {
			return item, i
		}
	}
	return nil, -1
}

func findByDescription(eas []*pb.EntityAnnotation, d string) (*pb.EntityAnnotation, int) {
	for i, item := range eas {
		if item.Description == d {
			return item, i
		}
	}
	return nil, -1
}

func joinDescriptions(eas []*pb.EntityAnnotation) string {
	var descriptions []string
	for _, ea := range eas {
		descriptions = append(descriptions, ea.Description)
	}
	return strings.Join(descriptions, " ")
}

func sortByAxis(eas []*pb.EntityAnnotation, axis string) {
	if eas != nil {
		sort.Slice(eas, func(i, j int) bool {
			return getField(eas[i].BoundingPoly.Vertices[3], axis) < getField(eas[j].BoundingPoly.Vertices[3], axis)
		})
	}
}

func anotherAxis(axis string) string {
	if axis == x {
		return y
	}
	return x
}

func printEAs(eas []*pb.EntityAnnotation) {
	for _, ea := range eas {
		log.Print(ea)
	}
	log.Print("--")
}
func findNearByInSameAxisAndSort(eas []*pb.EntityAnnotation, ea *pb.EntityAnnotation, axis string, avgDiff float64, sortAxis string) []*pb.EntityAnnotation {
	if ea == nil {
		return nil
	}
	var filtered = findNearByInSameAxis(eas, ea, axis, avgDiff)
	sortByAxis(filtered, sortAxis)
	return filtered
}

func findNearByInSameAxis(eas []*pb.EntityAnnotation, ea *pb.EntityAnnotation, axis string, avgDiff float64) []*pb.EntityAnnotation {
	if ea == nil {
		return nil
	}
	var filtered []*pb.EntityAnnotation

	if _, index := findByVertices(eas, ea.BoundingPoly.Vertices); index != -1 {
		// to the right
		for i := index; i < len(eas); i++ {
			if i == index { // first item
				filtered = append(filtered, eas[i])
			} else { // compares diff between this and the former point with the avg diff
				if float64(getField(eas[i].BoundingPoly.Vertices[3], axis)-
					getField(eas[i-1].BoundingPoly.Vertices[3], axis)) < avgDiff {
					filtered = append(filtered, eas[i])
				} else {
					break // considered no more nearby
				}
			}
		}
		// to the left
		for i := index - 1; i >= 0; i-- {
			if float64(getField(eas[i+1].BoundingPoly.Vertices[3], axis)-
				getField(eas[i].BoundingPoly.Vertices[3], axis)) < avgDiff {
				filtered = append(filtered, eas[i])
			} else {
				break // considered no more nearby
			}
		}
	}
	return filtered
}

const extraSameWordAVGSpace = 10

// diff of space between 2 words
func findNearBySameWord(eas []*pb.EntityAnnotation, ea *pb.EntityAnnotation, avgSpace int) []*pb.EntityAnnotation {
	if ea == nil {
		return nil
	}
	var filtered []*pb.EntityAnnotation

	if _, index := findByVertices(eas, ea.BoundingPoly.Vertices); index != -1 {
		// to the right
		for i := index; i < len(eas); i++ {
			if i == index { // first item
				filtered = append(filtered, eas[i])
			} else { // compares diff between this and the former point with the avg diff
				if int(eas[i].BoundingPoly.Vertices[3].X-
					eas[i-1].BoundingPoly.Vertices[2].X) <= avgSpace+extraSameWordAVGSpace {
					filtered = append(filtered, eas[i])
				} else {
					break // considered no more nearby
				}
			}
		}
		// to the left
		for i := index - 1; i >= 0; i-- {
			if int(eas[i+1].BoundingPoly.Vertices[3].X-
				eas[i].BoundingPoly.Vertices[2].X) <= avgSpace+extraSameWordAVGSpace {
				filtered = append(filtered, eas[i])
			} else {
				break // considered no more nearby
			}
		}
	}
	return filtered
}

func thirdPointDiffs(eas []*pb.EntityAnnotation, vertexField string) []int {
	var (
		out []int
	)
	if len(eas) == 0 {
		return []int{}
	}
	if len(eas) == 1 {
		return []int{getField(eas[0].BoundingPoly.Vertices[3], vertexField)}
	}
	for i := 1; i < len(eas); i++ {
		out = append(out,
			getField(eas[i].BoundingPoly.Vertices[3], vertexField)-
				getField(eas[i-1].BoundingPoly.Vertices[3], vertexField))
	}
	return out
}

func spaceDiffs(eas []*pb.EntityAnnotation) []int {
	var (
		out []int
	)
	if len(eas) < 1 {
		return []int{}
	}
	for i := 1; i < len(eas); i++ {
		out = append(out,
			int(eas[i].BoundingPoly.Vertices[3].X-
				eas[i-1].BoundingPoly.Vertices[2].X))
	}
	return out
}
