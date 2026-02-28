import { useEffect, useState } from 'react';
import { observer } from 'mobx-react-lite';
import { rootStore } from '../stores/RootStore';
import type { CompanyForm } from '../types';

export const Companies = observer(() => {
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [formData, setFormData] = useState<CompanyForm>({
    name: '',
    description: '',
  });

  useEffect(() => {
    rootStore.fetchCompanies();
  }, []);

  const { companies, loading, error } = rootStore;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingId) {
        await rootStore.updateCompany(editingId, formData);
      } else {
        await rootStore.createCompany(formData);
      }
      setFormData({ name: '', description: '' });
      setShowForm(false);
      setEditingId(null);
    } catch (err) {
      console.error(err);
    }
  };

  const handleEdit = (id: number) => {
    const company = companies.find((c) => c.id === id);
    if (company) {
      setFormData({ name: company.name, description: company.description });
      setEditingId(id);
      setShowForm(true);
    }
  };

  const handleDelete = async (id: number) => {
    if (confirm('Удалить компанию?')) {
      await rootStore.deleteCompany(id);
    }
  };

  const handleCancel = () => {
    setFormData({ name: '', description: '' });
    setShowForm(false);
    setEditingId(null);
  };

  if (loading.companies) {
    return <div className="loading-spinner">⏳ Загрузка...</div>;
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <h1 style={{ fontSize: '2rem', fontWeight: '700' }}>Компании</h1>
        <button className="btn btn-primary" onClick={() => setShowForm(!showForm)}>
          {showForm ? '❌ Закрыть' : '➕ Добавить компанию'}
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}

      {showForm && (
        <div className="card" style={{ marginBottom: '24px' }}>
          <h2 className="card-title">
            {editingId ? 'Редактировать компанию' : 'Новая компания'}
          </h2>
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label className="form-label">Название *</label>
              <input
                className="form-input"
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                required
              />
            </div>
            <div className="form-group">
              <label className="form-label">Описание</label>
              <textarea
                className="form-textarea"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              />
            </div>
            <div style={{ display: 'flex', gap: '12px' }}>
              <button type="submit" className="btn btn-primary">
                {editingId ? '💾 Сохранить' : '➕ Создать'}
              </button>
              <button type="button" className="btn btn-secondary" onClick={handleCancel}>
                ❌ Отмена
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="card">
        <table className="table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Название</th>
              <th>Описание</th>
              <th>Создано</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            {companies.map((company) => (
              <tr key={company.id}>
                <td>{company.id}</td>
                <td style={{ fontWeight: '600' }}>{company.name}</td>
                <td style={{ maxWidth: '400px' }}>
                  {company.description.substring(0, 100)}
                  {company.description.length > 100 && '...'}
                </td>
                <td>{new Date(company.created_at).toLocaleDateString('ru-RU')}</td>
                <td>
                  <div style={{ display: 'flex', gap: '8px' }}>
                    <button
                      className="btn btn-secondary"
                      onClick={() => handleEdit(company.id)}
                      style={{ padding: '6px 12px', fontSize: '12px' }}
                    >
                      ✏️ Изменить
                    </button>
                    <button
                      className="btn btn-danger"
                      onClick={() => handleDelete(company.id)}
                      style={{ padding: '6px 12px', fontSize: '12px' }}
                    >
                      🗑️ Удалить
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {companies.length === 0 && (
          <div style={{ textAlign: 'center', padding: '40px', color: '#6b7280' }}>
            Компаний пока нет
          </div>
        )}
      </div>
    </div>
  );
});
