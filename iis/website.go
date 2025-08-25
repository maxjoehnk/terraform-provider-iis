package iis

type Website struct {
	Name            string               `json:"name"`
	ID              string               `json:"id"`
	PhysicalPath    string               `json:"physical_path"`
	Bindings        []WebsiteBinding     `json:"bindings"`
	ApplicationPool ApplicationReference `json:"application_pool"`
}

type WebsiteBinding struct {
	Protocol  string `json:"protocol"`
	Port      int    `json:"port"`
	IPAddress string `json:"ip_address"`
	Hostname  string `json:"hostname"`
}
