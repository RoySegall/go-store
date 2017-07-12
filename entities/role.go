package entities

type Role struct {
	Title string `json,gorethink:"title"`
}

// Check if a role can do an operation on an entity.
func (role Role) AllowedAction(entity string, operation string) (bool) {
	return true
}
