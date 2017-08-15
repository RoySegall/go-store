package entities

type User struct {
	Id string `json,gorethink:"id,omitempty"`
	Username string `json,gorethink:"username"`
	Password string
	Email string `yml:"email"`
	Image string `yml:"Image"`
	Role Role `json,gorethink:"role"`
	Token Token
	Cart Cart `json,gorethink:"cart"`
	PastCarts []Cart `json,gorethink:"past_carts"`
}

type Token struct {
	Token string `json,gorethink:"token"`
	Expire int64 `json,gorethink:"expire"`
	RefreshToken string `json,gorethink:"refresh_token"`
}

type Cart struct {
	Items []Item `json,gorethink:"items"`
}

type PostItem struct {
	ItemId string `json,gorethink:"item_id"`
}
