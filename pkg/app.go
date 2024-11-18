package pkg

type Application struct {
	Env *Env
}

func App() (Application, error) {
	app := &Application{}
	app.Env = NewEnv()

	return *app, nil
}
