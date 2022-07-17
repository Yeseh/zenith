package init

import "errors"

func InitializerFactory(runtime string) (ZenithInitializer, error) {
	switch runtime {
	case "deno":
		return &DenoInitializer{}, nil
	}

	return nil, errors.New("Unknown runtime: " + runtime)
}
