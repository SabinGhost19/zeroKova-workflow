// Frontend - React application for test-workflow
// Version: 1.0.0
import React, { useState, useEffect } from 'react';
import { api } from './services/api';
import './App.css';

function App() {
  const [activeTab, setActiveTab] = useState('orders');
  const [orders, setOrders] = useState([]);
  const [products, setProducts] = useState([]);
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');

  // Form states
  const [orderForm, setOrderForm] = useState({
    customer_name: '',
    product_id: '',
    quantity: 1,
  });

  const [productForm, setProductForm] = useState({
    name: '',
    description: '',
    price: 0,
    quantity: 0,
  });

  useEffect(() => {
    loadData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeTab]);

  const loadData = async () => {
    setLoading(true);
    try {
      if (activeTab === 'orders') {
        const data = await api.getOrders();
        setOrders(data.orders || []);
      } else if (activeTab === 'inventory') {
        const data = await api.getProducts();
        setProducts(data.products || []);
      } else if (activeTab === 'notifications') {
        const data = await api.getNotifications();
        setNotifications(data.notifications || []);
      }
    } catch (error) {
      showMessage(`Error loading data: ${error.message}`, 'error');
    }
    setLoading(false);
  };

  const showMessage = (text, type = 'success') => {
    setMessage({ text, type });
    setTimeout(() => setMessage(''), 3000);
  };

  const handleCreateOrder = async () => {
    if (!orderForm.customer_name || !orderForm.product_id) {
      showMessage('Please fill all fields', 'error');
      return;
    }

    const selectedProduct = products.find(p => p.id === orderForm.product_id);
    if (!selectedProduct) {
      showMessage('Please select a valid product', 'error');
      return;
    }

    try {
      const result = await api.createOrder({
        customer_name: orderForm.customer_name,
        items: [{
          product_id: orderForm.product_id,
          product_name: selectedProduct.name,
          quantity: parseInt(orderForm.quantity),
          price: selectedProduct.price,
        }],
      });

      if (result.success) {
        showMessage('Order created successfully!');
        setOrderForm({ customer_name: '', product_id: '', quantity: 1 });
        loadData();
      } else {
        showMessage(result.message || 'Failed to create order', 'error');
      }
    } catch (error) {
      showMessage(`Error: ${error.message}`, 'error');
    }
  };

  const handleUpdateOrderStatus = async (orderId, status) => {
    try {
      const result = await api.updateOrderStatus(orderId, status);
      if (result.success) {
        showMessage(`Order status updated to ${status}`);
        loadData();
      } else {
        showMessage(result.message || 'Failed to update status', 'error');
      }
    } catch (error) {
      showMessage(`Error: ${error.message}`, 'error');
    }
  };

  const handleAddProduct = async () => {
    if (!productForm.name || productForm.price <= 0) {
      showMessage('Please fill all required fields', 'error');
      return;
    }

    try {
      const result = await api.addProduct(productForm);
      if (result.success) {
        showMessage('Product added successfully!');
        setProductForm({ name: '', description: '', price: 0, quantity: 0 });
        loadData();
      } else {
        showMessage(result.message || 'Failed to add product', 'error');
      }
    } catch (error) {
      showMessage(`Error: ${error.message}`, 'error');
    }
  };

  const handleUpdateStock = async (productId, change) => {
    try {
      const result = await api.updateStock(productId, change);
      if (result.success) {
        showMessage(`Stock updated by ${change > 0 ? '+' : ''}${change}`);
        loadData();
      } else {
        showMessage(result.message || 'Failed to update stock', 'error');
      }
    } catch (error) {
      showMessage(`Error: ${error.message}`, 'error');
    }
  };

  // Load products for order form
  useEffect(() => {
    api.getProducts().then(data => setProducts(data.products || []));
  }, []);

  return (
    <div className="app">
      <header className="header">
        <h1>Test Workflow - Zero Trust Demo</h1>
        <p>Orders & Inventory Microservices</p>
      </header>

      {message && (
        <div className={`message ${message.type}`}>
          {message.text}
        </div>
      )}

      <nav className="tabs">
        <button
          className={activeTab === 'orders' ? 'active' : ''}
          onClick={() => setActiveTab('orders')}
        >
          Orders
        </button>
        <button
          className={activeTab === 'inventory' ? 'active' : ''}
          onClick={() => setActiveTab('inventory')}
        >
          Inventory
        </button>
        <button
          className={activeTab === 'notifications' ? 'active' : ''}
          onClick={() => setActiveTab('notifications')}
        >
          Notifications
        </button>
      </nav>

      <main className="content">
        {loading && <div className="loading">Loading...</div>}

        {activeTab === 'orders' && (
          <div className="section">
            <h2>Create New Order</h2>
            <div className="form">
              <input
                type="text"
                placeholder="Customer Name"
                value={orderForm.customer_name}
                onChange={(e) => setOrderForm({ ...orderForm, customer_name: e.target.value })}
              />
              <select
                value={orderForm.product_id}
                onChange={(e) => setOrderForm({ ...orderForm, product_id: e.target.value })}
              >
                <option value="">Select Product</option>
                {products.map(p => (
                  <option key={p.id} value={p.id}>
                    {p.name} - ${p.price} (Stock: {p.quantity})
                  </option>
                ))}
              </select>
              <input
                type="number"
                min="1"
                placeholder="Quantity"
                value={orderForm.quantity}
                onChange={(e) => setOrderForm({ ...orderForm, quantity: e.target.value })}
              />
              <button onClick={handleCreateOrder} className="btn-primary">
                Create Order (REST → gRPC → Kafka)
              </button>
            </div>

            <h2>Orders List</h2>
            <button onClick={loadData} className="btn-secondary">Refresh</button>
            <div className="list">
              {orders.length === 0 ? (
                <p>No orders yet</p>
              ) : (
                orders.map(order => (
                  <div key={order.id} className="card">
                    <div className="card-header">
                      <strong>Order: {order.id?.substring(0, 8)}...</strong>
                      <span className={`status ${order.status?.toLowerCase()}`}>
                        {order.status}
                      </span>
                    </div>
                    <div className="card-body">
                      <p>Customer: {order.customer_name}</p>
                      <p>Total: ${order.total_amount?.toFixed(2)}</p>
                      <p>Items: {order.items?.length || 0}</p>
                    </div>
                    <div className="card-actions">
                      {order.status === 'PENDING' && (
                        <>
                          <button onClick={() => handleUpdateOrderStatus(order.id, 'CONFIRMED')} className="btn-success">
                            Confirm
                          </button>
                          <button onClick={() => handleUpdateOrderStatus(order.id, 'CANCELLED')} className="btn-danger">
                            Cancel
                          </button>
                        </>
                      )}
                      {order.status === 'CONFIRMED' && (
                        <button onClick={() => handleUpdateOrderStatus(order.id, 'SHIPPED')} className="btn-primary">
                          Ship
                        </button>
                      )}
                      {order.status === 'SHIPPED' && (
                        <button onClick={() => handleUpdateOrderStatus(order.id, 'DELIVERED')} className="btn-success">
                          Delivered
                        </button>
                      )}
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        )}

        {activeTab === 'inventory' && (
          <div className="section">
            <h2>Add New Product</h2>
            <div className="form">
              <input
                type="text"
                placeholder="Product Name"
                value={productForm.name}
                onChange={(e) => setProductForm({ ...productForm, name: e.target.value })}
              />
              <input
                type="text"
                placeholder="Description"
                value={productForm.description}
                onChange={(e) => setProductForm({ ...productForm, description: e.target.value })}
              />
              <input
                type="number"
                min="0"
                step="0.01"
                placeholder="Price"
                value={productForm.price}
                onChange={(e) => setProductForm({ ...productForm, price: parseFloat(e.target.value) })}
              />
              <input
                type="number"
                min="0"
                placeholder="Initial Quantity"
                value={productForm.quantity}
                onChange={(e) => setProductForm({ ...productForm, quantity: parseInt(e.target.value) })}
              />
              <button onClick={handleAddProduct} className="btn-primary">
                Add Product (REST → gRPC)
              </button>
            </div>

            <h2>Products List</h2>
            <button onClick={loadData} className="btn-secondary">Refresh</button>
            <div className="list">
              {products.length === 0 ? (
                <p>No products yet</p>
              ) : (
                products.map(product => (
                  <div key={product.id} className="card">
                    <div className="card-header">
                      <strong>{product.name}</strong>
                      <span className="price">${product.price?.toFixed(2)}</span>
                    </div>
                    <div className="card-body">
                      <p>{product.description}</p>
                      <p>Stock: {product.quantity} | Reserved: {product.reserved || 0}</p>
                    </div>
                    <div className="card-actions">
                      <button onClick={() => handleUpdateStock(product.id, 10)} className="btn-success">
                        +10 Stock
                      </button>
                      <button onClick={() => handleUpdateStock(product.id, -5)} className="btn-danger">
                        -5 Stock
                      </button>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        )}

        {activeTab === 'notifications' && (
          <div className="section">
            <h2>Notifications</h2>
            <button onClick={loadData} className="btn-secondary">Refresh</button>
            <div className="list">
              {notifications.length === 0 ? (
                <p>No notifications yet</p>
              ) : (
                notifications.map(notification => (
                  <div key={notification.id} className="card">
                    <div className="card-header">
                      <span className={`type ${notification.type?.toLowerCase()}`}>
                        {notification.type}
                      </span>
                      <small>{notification.created_at}</small>
                    </div>
                    <div className="card-body">
                      <p>{notification.message}</p>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        )}
      </main>

      <footer className="footer">
        <p>Microservices: Go (API Gateway) → Java (Orders) → C# (Inventory) → Ruby (Notifications)</p>
        <p>Communication: REST + gRPC + Kafka | Database: PostgreSQL</p>
      </footer>
    </div>
  );
}

export default App;
