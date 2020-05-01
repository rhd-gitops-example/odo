package webhook

// Create create a new webhook on the target Git Repository
func Create(accessToken, pipelines, serviceName string, isCICD, isInsecure bool) error {
	return nil
}

// Delete deletes webhooks on the target Git Repository that match the listener address
func Delete(accessToken, pipelines, serviceName string, isCICD, isInsecure bool) error {
	return nil
}
