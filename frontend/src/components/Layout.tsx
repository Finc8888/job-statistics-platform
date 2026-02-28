import { ReactNode } from 'react';
import { Link, useLocation } from 'react-router-dom';

interface LayoutProps {
  children: ReactNode;
}

export const Layout = ({ children }: LayoutProps) => {
  const location = useLocation();

  const navItems = [
    { path: '/', label: 'Дашборд' },
    { path: '/companies', label: 'Компании' },
    { path: '/jobs', label: 'Вакансии' },
    { path: '/skills', label: 'Навыки' },
    { path: '/statistics', label: 'Статистика' },
  ];

  return (
    <div style={{ minHeight: '100vh', background: '#f5f7fa' }}>
      <header
        style={{
          background: 'white',
          boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
          position: 'sticky',
          top: 0,
          zIndex: 100,
        }}
      >
        <div
          className="container"
          style={{
            display: 'flex',
            alignItems: 'center',
            height: '64px',
            gap: '32px',
          }}
        >
          <h1 style={{ fontSize: '1.5rem', fontWeight: '700', color: '#3b82f6' }}>
            📊 Job Stats
          </h1>
          <nav style={{ display: 'flex', gap: '24px' }}>
            {navItems.map((item) => (
              <Link
                key={item.path}
                to={item.path}
                style={{
                  textDecoration: 'none',
                  color: location.pathname === item.path ? '#3b82f6' : '#6b7280',
                  fontWeight: location.pathname === item.path ? '600' : '500',
                  fontSize: '14px',
                  transition: 'color 0.2s',
                }}
              >
                {item.label}
              </Link>
            ))}
          </nav>
        </div>
      </header>
      <main className="container" style={{ paddingTop: '32px', paddingBottom: '32px' }}>
        {children}
      </main>
    </div>
  );
};
