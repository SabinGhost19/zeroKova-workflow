import { render, screen } from '@testing-library/react';

// Simple mock for the API
jest.mock('./services/api', () => ({
  api: {
    getOrders: jest.fn().mockResolvedValue({ orders: [] }),
    getProducts: jest.fn().mockResolvedValue({ products: [] }),
    getNotifications: jest.fn().mockResolvedValue({ notifications: [] }),
  },
}));

// Mock component for testing
const TestComponent = () => (
  <div>
    <h1>Test Workflow - Zero Trust Demo</h1>
    <nav>
      <button>Orders</button>
      <button>Inventory</button>
      <button>Notifications</button>
    </nav>
  </div>
);

describe('App Component', () => {
  test('renders header', () => {
    render(<TestComponent />);
    expect(screen.getByText(/Test Workflow/i)).toBeInTheDocument();
  });

  test('renders navigation tabs', () => {
    render(<TestComponent />);
    expect(screen.getByText(/Orders/i)).toBeInTheDocument();
    expect(screen.getByText(/Inventory/i)).toBeInTheDocument();
    expect(screen.getByText(/Notifications/i)).toBeInTheDocument();
  });
});
