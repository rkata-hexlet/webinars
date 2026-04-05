package main

// Группы позволяют организовать роуты логически и применять middleware к группе.

/*
func groups() {
	r := gin.Default()
	api := r.Group("/api")
	api.Use(AuthMiddleware()) // Применяется ко всем роутам в группе
	{
		api.GET("/users", getUsers)    // /api/users
		api.POST("/users", createUser) // /api/users
		api.GET("/posts", getPosts)    // /api/posts
	}

	api1 := r.Group("/api")
	{
		v1 := api1.Group("/v1")
		{
			v1.GET("/users", getUsersV1) // /api/v1/users
		}

		v2 := api.Group("/v2")
		{
			v2.GET("/users", getUsersV2) // /api/v2/users
		}
	}
}
*/
