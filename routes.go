type Route struct {
        Method string
        Path   string
        Handle mux.Handle   // httprouter package as mux
}

type Routes []Route

var routes = Routes{
        Route{
                "GET",
                "/",
                Index,
        },
        Route{
                "GET",
                "/posts"
                PostIndex,
        },
        Route{
                "GET",
                "/posts/:id",
                PostShow,
        },
        Route{
                "POST",
                "/posts",
                PostCreate,
        },
}