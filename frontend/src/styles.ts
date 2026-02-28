export const colors = {
  primary: '#3b82f6',
  primaryDark: '#2563eb',
  secondary: '#8b5cf6',
  success: '#10b981',
  warning: '#f59e0b',
  danger: '#ef4444',
  gray: {
    50: '#f9fafb',
    100: '#f3f4f6',
    200: '#e5e7eb',
    300: '#d1d5db',
    400: '#9ca3af',
    500: '#6b7280',
    600: '#4b5563',
    700: '#374151',
    800: '#1f2937',
    900: '#111827',
  },
};

export const styles = `
  .container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 0 20px;
  }

  .card {
    background: white;
    border-radius: 12px;
    padding: 24px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    margin-bottom: 20px;
  }

  .card-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1f2937;
    margin-bottom: 16px;
  }

  .btn {
    padding: 10px 20px;
    border-radius: 8px;
    border: none;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 14px;
  }

  .btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }

  .btn-primary {
    background: ${colors.primary};
    color: white;
  }

  .btn-primary:hover {
    background: ${colors.primaryDark};
  }

  .btn-secondary {
    background: ${colors.gray[200]};
    color: ${colors.gray[700]};
  }

  .btn-danger {
    background: ${colors.danger};
    color: white;
  }

  .btn-success {
    background: ${colors.success};
    color: white;
  }

  .table {
    width: 100%;
    border-collapse: collapse;
  }

  .table th,
  .table td {
    padding: 12px;
    text-align: left;
    border-bottom: 1px solid ${colors.gray[200]};
  }

  .table th {
    background: ${colors.gray[50]};
    font-weight: 600;
    color: ${colors.gray[700]};
  }

  .table tr:hover {
    background: ${colors.gray[50]};
  }

  .form-group {
    margin-bottom: 16px;
  }

  .form-label {
    display: block;
    margin-bottom: 8px;
    font-weight: 500;
    color: ${colors.gray[700]};
  }

  .form-input,
  .form-select,
  .form-textarea {
    width: 100%;
    padding: 10px 12px;
    border: 1px solid ${colors.gray[300]};
    border-radius: 8px;
    font-size: 14px;
    transition: border-color 0.2s;
  }

  .form-input:focus,
  .form-select:focus,
  .form-textarea:focus {
    outline: none;
    border-color: ${colors.primary};
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }

  .form-textarea {
    resize: vertical;
    min-height: 100px;
  }

  .form-checkbox {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .badge {
    display: inline-block;
    padding: 4px 12px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 500;
  }

  .badge-primary {
    background: rgba(59, 130, 246, 0.1);
    color: ${colors.primary};
  }

  .badge-success {
    background: rgba(16, 185, 129, 0.1);
    color: ${colors.success};
  }

  .badge-warning {
    background: rgba(245, 158, 11, 0.1);
    color: ${colors.warning};
  }

  .loading-spinner {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 40px;
  }

  .error-message {
    background: rgba(239, 68, 68, 0.1);
    color: ${colors.danger};
    padding: 12px 16px;
    border-radius: 8px;
    margin-bottom: 16px;
  }

  .grid {
    display: grid;
    gap: 20px;
  }

  .grid-2 {
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  }

  .stat-card {
    background: white;
    border-radius: 12px;
    padding: 20px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .stat-value {
    font-size: 2rem;
    font-weight: 700;
    color: ${colors.gray[900]};
  }

  .stat-label {
    font-size: 0.875rem;
    color: ${colors.gray[600]};
    margin-top: 4px;
  }
`;
