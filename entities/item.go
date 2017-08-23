package entities

type Item struct {
	Id 					string 	`gorethink:"id,omitempty"`
	Title 			string 	`json,gorethink:"title" form:"title"`
	Description string 	`json,gorethink:"description" form:"description"`
	Price 			float64 `json,gorethink:"price" form:"price"`
	Image 			string 	`json,gorethink:"image" form:"image"`
}
