package entities

type Item struct {
	Id 		string 	`json,gorethink:"id,Id"`
	Title string 	`json,gorethink:"title"`
	Price float64 `json,gorethink:"price"`
	Image string 	`json,gorethink:"image"`
}
