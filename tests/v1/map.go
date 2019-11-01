package v1

// TestGetGameMap test returning json object when valid paramters are given
func (t *APITest) TestGetGameMap() {
	t.Get("/api/v1/gamemap?cx=0&cy=0&scale=8&delegate=4")
	t.AssertOk()
}

// TestGetGameMapInvalid test returning 422 status code when invalid parameters are given
func (t *APITest) TestGetGameMapInvalid() {
	t.Get("/api/v1/gamemap")
	t.AssertStatus(422)
}
