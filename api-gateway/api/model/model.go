package model

type CreateUserRequestBody struct {
	Id        string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	FirstName string  `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name"`
	LastName  string  `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name"`
	Posts     []*Post `protobuf:"bytes,4,rep,name=posts,proto3" json:"posts"`
}

type Post struct {
	Id          string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Name        string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	Description string  `protobuf:"bytes,3,opt,name=description,proto3" json:"description"`
	UserId      string  `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id"`
	Medias      []*Post `protobuf:"bytes,5,rep,name=medias,proto3" json:"medias"`
}

type RegisterModel struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Code      int64  `json:"code"`
}

type RegisterResponseModel struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

type JwtRequestModel struct {
	Token string `json:"token"`
}


type ResponseError struct{
	Error interface{} `json:"error"`
}

// ServerError ...
type ServerError struct{
	Status string `json:"status"`
	Message string`json:"message"`
}
