package appy

type nilWebsocketFactory struct{}

type nilWebsocket struct{}

func (n *nilWebsocketFactory) Initialize(app *Appy, options WebsocketFactoryOptions) error {
	return nil
}

func (n *nilWebsocketFactory) Create(options WebsocketOptions) Websocket {
	return &nilWebsocket{}
}

func (n *nilWebsocket) Spin(c *HttpContext) error {
	return nil
}

func (n *nilWebsocket) Send(data []byte) {
}

func (n *nilWebsocket) Close() error {
	return nil
}
