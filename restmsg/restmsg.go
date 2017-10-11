package restmsg

type StoreRequest struct {
	File []byte
}

type StoreResponse struct {
	Status int
	Message string
	ID string
}

type CatResponse struct {
	Status int
	Message string
	File []byte
}

type GenericResponse struct {
	Status int
	Message string
}
