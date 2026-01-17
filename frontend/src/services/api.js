const API_BASE = process.env.REACT_APP_API_URL || '/api/v1';

export const api = {
  // Orders
  async createOrder(orderData) {
    const response = await fetch(`${API_BASE}/orders`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(orderData),
    });
    return response.json();
  },

  async getOrders(limit = 10, offset = 0) {
    const response = await fetch(`${API_BASE}/orders?limit=${limit}&offset=${offset}`);
    return response.json();
  },

  async getOrder(id) {
    const response = await fetch(`${API_BASE}/orders/${id}`);
    return response.json();
  },

  async updateOrderStatus(id, status) {
    const response = await fetch(`${API_BASE}/orders/${id}/status`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ status }),
    });
    return response.json();
  },

  // Inventory
  async getProducts(limit = 10, offset = 0) {
    const response = await fetch(`${API_BASE}/inventory/products?limit=${limit}&offset=${offset}`);
    return response.json();
  },

  async addProduct(productData) {
    const response = await fetch(`${API_BASE}/inventory/products`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(productData),
    });
    return response.json();
  },

  async updateStock(productId, quantityChange) {
    const response = await fetch(`${API_BASE}/inventory/products/${productId}/stock`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ quantity_change: quantityChange }),
    });
    return response.json();
  },

  // Notifications
  async getNotifications(limit = 10, offset = 0) {
    const response = await fetch(`${API_BASE}/notifications?limit=${limit}&offset=${offset}`);
    return response.json();
  },
};
