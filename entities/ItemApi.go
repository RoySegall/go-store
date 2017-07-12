package entities

import "store/api"

// Insert an object.
func (item Item) Insert() (string) {
	return api.Insert("item", item)
}

// Get a single object.
func (item Item) Get(id string) (Item) {
	res := api.Get("item", id)
	items := []Item{}
	res.All(&items)
	return items[0]
}

// Get all the items.
func (item Item) GetAll() ([]Item) {
	res := api.GetAll("item")
	items := []Item{}
	res.All(&items)
	return items
}

// Delete an item.
func (item Item) Delete() {
	api.Delete("item", item.Id)
}

// Update an item.
func (item Item) Update() {
	api.Update("item", item)
}
