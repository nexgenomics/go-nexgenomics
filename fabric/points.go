package fabric


type ID string

// Point
type Point struct {
	Id: ID
	Embedding: []float32
	Content: []byte
	Metadata: map[string]string
	// scores, etc
}


// UpsertPoint
func UpsertPoint(pt *Point) (ID,error) {


}

