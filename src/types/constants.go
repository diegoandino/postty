package types

// HTTPMethods contains all supported HTTP methods
var HTTPMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

// ContentTypes contains all supported content types
var ContentTypes = []string{
	"application/json",
	"application/xml",
	"text/plain",
	"application/x-www-form-urlencoded",
	"multipart/form-data",
}

// HeaderTemplates contains predefined header templates
var HeaderTemplates = []HeaderTemplate{
	{Name: "Authorization (Bearer)", Key: "Authorization", Placeholder: "Bearer <your-token>"},
	{Name: "Authorization (Basic)", Key: "Authorization", Placeholder: "Basic <base64-credentials>"},
	{Name: "API Key", Key: "X-API-Key", Placeholder: "<your-api-key>"},
	{Name: "Cookie", Key: "Cookie", Placeholder: "session_id=<value>"},
	{Name: "User Agent", Key: "User-Agent", Placeholder: "MyApp/1.0"},
	{Name: "Accept", Key: "Accept", Placeholder: "application/json"},
	{Name: "Custom Header", Key: "", Placeholder: ""},
}
