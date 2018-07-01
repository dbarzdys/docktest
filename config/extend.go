package config

func extend(to, from Config) Config {
	to.Export = extendStringMap(
		to.Export,
		from.Export,
	)
	to.Constants = extendStringMap(
		to.Constants,
		from.Constants,
	)
	to.Services = extendServices(
		to.Services,
		from.Services,
	)
	return to
}

func extendStringMap(to, from map[string]string) map[string]string {
	if to == nil {
		to = make(map[string]string)
	}
	if from == nil {
		from = make(map[string]string)
	}
	for k, v := range from {
		_, ok := to[k]
		if !ok {
			to[k] = v
		}
	}
	return to
}

func extendStringSlice(to, from []string) []string {
Loop:
	for _, f := range from {
		for _, t := range to {
			if t == f {
				continue Loop
			}
		}
		to = append(to, f)
	}
	return to
}

func extendString(to, from string) string {
	if to == "" {
		return from
	}
	return to
}

func extendServices(to, from map[string]Service) map[string]Service {
	if to == nil {
		to = make(map[string]Service)
	}
	if from == nil {
		from = make(map[string]Service)
	}
	for k, v := range from {
		found, ok := to[k]
		if !ok {
			found = v
		} else {
			found.Env = extendStringMap(found.Env, v.Env)
			found.Image = extendString(found.Image, v.Image)
			found.Tag = extendString(found.Tag, v.Tag)
			found.DependsOn = extendStringSlice(found.DependsOn, v.DependsOn)
			found.HealthCheck = extendStringSlice(found.HealthCheck, v.HealthCheck)
		}
		to[k] = found
	}
	return to
}
