package test

deployment: {
	// Multiple nested resources
	test: {
		deployment: {
			apiVersion: "apps/v1"
			kind:       "Deployment"
			metadata: {
				name:      "myapp"
				namespace: "mynamespace"
			}
		}
		service: {
			apiVersion: ""
			kind:       "Service"
			metadata: {
				name:      "myservice"
				namespace: "mynamespace"
			}
		}
	}
}
// No namespace test
prod: {
	apiVersion: "apps/v1"
	kind:       "Deployment"
	metadata: {
		name: "myapp"
	}
}
