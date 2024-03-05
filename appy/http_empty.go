package appy

type nilHttpProvider struct {
}

type nilHttpEndpointGroup struct {
}

func (n *nilHttpProvider) Initialize(*Appy, HttpOptions) error {
	return nil
}

func (n *nilHttpProvider) Run() error {
	return nil
}

func (n *nilHttpProvider) RootGroup() HttpEndpointGroup {
	return &nilHttpEndpointGroup{}
}

func (n *nilHttpEndpointGroup) Subgroup(string) HttpEndpointGroup {
	return &nilHttpEndpointGroup{}
}

func (n *nilHttpEndpointGroup) Use(HttpMiddleware) {
}

func (n *nilHttpEndpointGroup) GET(string, HttpHandler) {
}

func (n *nilHttpEndpointGroup) POST(string, HttpHandler) {
}

func (n *nilHttpEndpointGroup) PUT(string, HttpHandler) {
}

func (n *nilHttpEndpointGroup) DELETE(string, HttpHandler) {
}
