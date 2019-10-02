package cloudvision

import (
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

const x = "X"
const y = "Y"

type Document struct {
	document        []*pb.EntityAnnotation
	hSortedDocument []*pb.EntityAnnotation // sort by x axis
	vSortedDocument []*pb.EntityAnnotation // sort by y axis
	hAVGDiff        float64
	vAVGDiff        float64
}

func NewDocument(d []*pb.EntityAnnotation) Document {
	doc := Document{document: d}
	doc.hSortedDocument = doc.cloneDocument(1) // excludes the first
	doc.vSortedDocument = doc.cloneDocument(1) // excludes the first
	doc.sortByThirdPoint(x)
	doc.sortByThirdPoint(y)
	doc.hAVGDiff = average(doc.thirdPointDiffs(x))
	doc.vAVGDiff = average(doc.thirdPointDiffs(y))
	return doc
}

func (d *Document) Parse() PassportInfo {
	var pi = PassportInfo{}

	pi.Type = d.extractByLabel("Type")
	pi.Name = d.extractByLabelAndJoin("Name")
	pi.Sex = d.extractByLabel("Sex")
	pi.CountryCode = d.extractByLabel("Country")
	pi.PassportNo = d.extractByLabelAndJoin("Passport")
	pi.Nationality = d.extractByLabel("Nationality")
	pi.DateOfBirth = d.extractByLabelContainsAndJoin("birth")
	pi.DateOfIssue = d.extractByLabelContainsAndJoin("issue")
	pi.DateOfExpiry = d.extractByLabelContainsAndJoin("expiry")

	return pi
}

func (d *Document) extractByLabel(label string) string {
	ea := d.findByDescriptions(label)
	if ea == nil {
		return ""
	}

	nearby := d.findNearByInSameAxis(ea, x)
	// printEAs(nearby)
	_, index := findByDescription(nearby, label)
	if index == -1 {
		return ""
	}
	valEA := nearby[index+1]
	return valEA.Description
}

func (d *Document) extractByLabelAndJoin(label string) string {
	ea := d.findByDescriptions(label)
	if ea == nil {
		return ""
	}
	nearby := d.findNearByInSameAxis(ea, x)
	// printEAs(nearby)
	_, index := findByDescription(nearby, label)
	if index == -1 {
		return ""
	}
	valEA := nearby[index+1]
	valPlane := d.findNearByInSameAxis(valEA, y)
	filtered := findNearBySameWord(valPlane, valEA, int(average(spaceDiffs(valPlane))))
	sortByAxis(filtered, x)
	// printEAs(filtered)
	return joinDescriptions(filtered)
}

func (d *Document) extractByLabelContainsAndJoin(label string) string {
	var labelHeadEA = d.findLabelHead(label)
	if labelHeadEA != nil {
		nearby := d.findNearByInSameAxis(labelHeadEA, x)
		// printEAs(nearby)
		_, index := findByVertices(nearby, labelHeadEA.BoundingPoly.Vertices)
		if index != -1 {
			valEA := nearby[index+1]
			valPlane := d.findNearByInSameAxis(valEA, y)
			avg := int(average(spaceDiffs(valPlane)))
			// log.Print(avg)
			filtered := findNearBySameWord(valPlane, valEA, avg)
			// printEAs(filtered)

			sortByAxis(filtered, x)
			return joinDescriptions(filtered)
		}
	}

	return ""
}

func (d *Document) findLabelHead(label string) *pb.EntityAnnotation {
	ea := d.findByDescriptions(label)
	var labelHeadEA *pb.EntityAnnotation
	if ea != nil {
		nearby := d.findNearByInSameAxis(ea, y)
		avg := int(average(spaceDiffs(nearby)))
		// log.Print(avg)
		filtered := findNearBySameWord(nearby, ea, avg)
		// printEAs(filtered)
		sortByAxis(filtered, x)

		labelHeadEA = filtered[0]
	}
	return labelHeadEA
}

func (d *Document) findByDescriptions(values ...string) *pb.EntityAnnotation {
	for _, v := range values {
		for _, ea := range d.document {
			if ea.Description == v {
				return ea
			}
		}
	}
	return nil
}

func (d *Document) findNearByInSameAxis(ea *pb.EntityAnnotation, axis string) []*pb.EntityAnnotation {
	if ea == nil {
		return nil
	}
	return findNearByInSameAxisAndSort(
		d.getSortedDocuments(axis), ea, axis, d.getAVGDiff(axis), anotherAxis(axis))
}

func (d *Document) findStartToEndAxis(ea *pb.EntityAnnotation, axis string) []*pb.EntityAnnotation {
	return nil
}

func (d *Document) cloneDocument(firstIndex int) []*pb.EntityAnnotation {
	return append([]*pb.EntityAnnotation(nil), d.document[firstIndex:]...)
}

func (d *Document) sortByThirdPoint(vertexField string) {
	sortByAxis(d.getSortedDocuments(vertexField), vertexField)
}

func (d *Document) thirdPointDiffs(vertexField string) []int {
	return thirdPointDiffs(d.getSortedDocuments(vertexField), vertexField)
}

func (d *Document) getSortedDocuments(vertexField string) []*pb.EntityAnnotation {
	var doc []*pb.EntityAnnotation
	switch vertexField {
	case x:
		doc = d.hSortedDocument
	case y:
		doc = d.vSortedDocument
	}
	return doc
}

func (d *Document) getAVGDiff(vertexField string) float64 {
	var avgDiff float64
	switch vertexField {
	case x:
		avgDiff = d.hAVGDiff
	case y:
		avgDiff = d.vAVGDiff
	}
	return avgDiff
}
