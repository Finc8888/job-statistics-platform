import { useEffect, useState } from 'react';
import { observer } from 'mobx-react-lite';
import { rootStore } from '../stores/RootStore';
import type { SkillForm } from '../types';

export const Skills = observer(() => {
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [formData, setFormData] = useState<SkillForm>({
    name: '',
    category: 'Язык программирования',
  });
  const [filterCategory, setFilterCategory] = useState<string>('all');

  useEffect(() => {
    rootStore.fetchSkills();
  }, []);

  const { skills, loading, error } = rootStore;

  const categories = [
    'Язык программирования',
    'База данных',
    'Фреймворк',
    'Инструмент',
    'Другое',
  ];

  const filteredSkills = filterCategory === 'all'
    ? skills
    : skills.filter((s) => s.category === filterCategory);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingId) {
        await rootStore.updateSkill(editingId, formData);
      } else {
        await rootStore.createSkill(formData);
      }
      setFormData({ name: '', category: 'Язык программирования' });
      setShowForm(false);
      setEditingId(null);
    } catch (err) {
      console.error(err);
    }
  };

  const handleEdit = (id: number) => {
    const skill = skills.find((s) => s.id === id);
    if (skill) {
      setFormData({ name: skill.name, category: skill.category });
      setEditingId(id);
      setShowForm(true);
    }
  };

  const handleDelete = async (id: number) => {
    if (confirm('Удалить навык?')) {
      await rootStore.deleteSkill(id);
    }
  };

  const getCategoryColor = (category: string) => {
    const colors: Record<string, string> = {
      'Язык программирования': 'primary',
      'База данных': 'success',
      'Фреймворк': 'warning',
      'Инструмент': 'badge-primary',
      'Другое': 'badge-secondary',
    };
    return colors[category] || 'primary';
  };

  if (loading.skills) {
    return <div className="loading-spinner">⏳ Загрузка...</div>;
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <h1 style={{ fontSize: '2rem', fontWeight: '700' }}>Навыки</h1>
        <button className="btn btn-primary" onClick={() => setShowForm(!showForm)}>
          {showForm ? '❌ Закрыть' : '➕ Добавить навык'}
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}

      {/* Filter */}
      <div className="card" style={{ marginBottom: '24px' }}>
        <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap' }}>
          <button
            className={`btn ${filterCategory === 'all' ? 'btn-primary' : 'btn-secondary'}`}
            onClick={() => setFilterCategory('all')}
          >
            Все ({skills.length})
          </button>
          {categories.map((cat) => (
            <button
              key={cat}
              className={`btn ${filterCategory === cat ? 'btn-primary' : 'btn-secondary'}`}
              onClick={() => setFilterCategory(cat)}
            >
              {cat} ({skills.filter((s) => s.category === cat).length})
            </button>
          ))}
        </div>
      </div>

      {showForm && (
        <div className="card" style={{ marginBottom: '24px' }}>
          <h2 className="card-title">
            {editingId ? 'Редактировать навык' : 'Новый навык'}
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
              <label className="form-label">Категория *</label>
              <select
                className="form-select"
                value={formData.category}
                onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                required
              >
                {categories.map((cat) => (
                  <option key={cat} value={cat}>
                    {cat}
                  </option>
                ))}
              </select>
            </div>
            <div style={{ display: 'flex', gap: '12px' }}>
              <button type="submit" className="btn btn-primary">
                {editingId ? '💾 Сохранить' : '➕ Создать'}
              </button>
              <button
                type="button"
                className="btn btn-secondary"
                onClick={() => {
                  setFormData({ name: '', category: 'Язык программирования' });
                  setShowForm(false);
                  setEditingId(null);
                }}
              >
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
              <th>Категория</th>
              <th>Создано</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            {filteredSkills.map((skill) => (
              <tr key={skill.id}>
                <td>{skill.id}</td>
                <td style={{ fontWeight: '600' }}>{skill.name}</td>
                <td>
                  <span className={`badge ${getCategoryColor(skill.category)}`}>
                    {skill.category}
                  </span>
                </td>
                <td>{new Date(skill.created_at).toLocaleDateString('ru-RU')}</td>
                <td>
                  <div style={{ display: 'flex', gap: '8px' }}>
                    <button
                      className="btn btn-secondary"
                      onClick={() => handleEdit(skill.id)}
                      style={{ padding: '6px 12px', fontSize: '12px' }}
                    >
                      ✏️ Изменить
                    </button>
                    <button
                      className="btn btn-danger"
                      onClick={() => handleDelete(skill.id)}
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
        {filteredSkills.length === 0 && (
          <div style={{ textAlign: 'center', padding: '40px', color: '#6b7280' }}>
            Навыков в этой категории пока нет
          </div>
        )}
      </div>
    </div>
  );
});
