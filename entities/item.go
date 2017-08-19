package entities

type Item struct {
	Id 					string 	`gorethink:"id,omitempty"`
	Title 			string 	`json,gorethink:"title"`
	Description string 	`json,gorethink:"description"`
	Price 			float64 `json,gorethink:"price"`
	Image 			string 	`json,gorethink:"image"`
}
