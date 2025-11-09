class WebSocketClient {
  constructor(url, options = {}) {
    this.url = url;
    this.options = {
      reconnectInterval: 1000,
      maxReconnectAttempts: 5,
      heartbeatInterval: 30000,
      ...options,
    };
    this.reconnectAttempts = 0;
    this.messageQueue = [];
    this.eventHandlers = {};
    this.isConnected = false;

    this.connect();
  }

  connect() {
    console.log(`Connecting to ${this.url}...`);

    try {
      this.ws = new WebSocket(this.url);
      this.setupEventHandlers();
    } catch (error) {
      console.error('Failed to create WebSocket:', error);
      this.scheduleReconnect();
    }
  }

  send(message) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      const data =
        typeof message === 'object' ? JSON.stringify(message) : message;
      this.ws.send(data);
    } else {
      console.log("WebSocket not connected");
    }
  }

  setupEventHandlers() {
    this.ws.onopen = (event) => {
      console.log('WebSocket connected');
      this.isConnected = true;
      this.reconnectAttempts = 0;
    };
  }

  scheduleReconnect() {
    if (this.reconnectAttempts < this.options.maxReconnectAttempts) {
      this.reconnectAttempts++;
      console.log(`Reconnect attempt ${this.reconnectAttempts}...`);
      setTimeout(() => this.connect(), this.options.reconnectInterval);
    } else {
      console.warn('Max reconnect attempts reached. Giving up.');
    }
  }
}