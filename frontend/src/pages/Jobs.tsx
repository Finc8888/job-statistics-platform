import { useEffect, useState } from 'react';
import { observer } from 'mobx-react-lite';
import { rootStore } from '../stores/RootStore';
import type { JobForm } from '../types';

export const Jobs = observer(() => {
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [selectedSkillIds, setSelectedSkillIds] = useState<number[]>([]);
  const [formData, setFormData] = useState<JobForm>({
    company_id: 0,
    title: '',
    level: 'Middle',
    specialization: '',
    salary_min: undefined,
    salary_max: undefined,
    salary_currency: 'RUB',
    experience_years: '',
    location: '',
    remote_available: false,
    description: '',
    responsibilities: '',
    benefits: '',
    is_active: true,
  });

  useEffect(() => {
    const loadJobs = async () => {
      await rootStore.fetchJobs();
      rootStore.fetchAllJobSkills();
    };
    loadJobs();
    rootStore.fetchCompanies();
    rootStore.fetchSkills();
  }, []);

  const { jobs, companies, skills, loading, error } = rootStore;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingId) {
        await rootStore.updateJob(editingId, formData);
        await rootStore.setJobSkills(editingId, selectedSkillIds);
      } else {
        const job = await rootStore.createJob(formData);
        await rootStore.setJobSkills(job.id, selectedSkillIds);
      }
      resetForm();
    } catch (err) {
      console.error(err);
    }
  };

  const resetForm = () => {
    setFormData({
      company_id: 0,
      title: '',
      level: 'Middle',
      specialization: '',
      salary_min: undefined,
      salary_max: undefined,
      salary_currency: 'RUB',
      experience_years: '',
      location: '',
      remote_available: false,
      description: '',
      responsibilities: '',
      benefits: '',
      is_active: true,
    });
    setSelectedSkillIds([]);
    setShowForm(false);
    setEditingId(null);
  };

  const handleEdit = async (id: number) => {
    const job = jobs.find((j) => j.id === id);
    if (job) {
      setFormData({
        company_id: job.company_id,
        title: job.title,
        level: job.level,
        specialization: job.specialization || '',
        salary_min: job.salary_min || undefined,
        salary_max: job.salary_max || undefined,
        salary_currency: job.salary_currency,
        experience_years: job.experience_years || '',
        location: job.location || '',
        remote_available: job.remote_available,
        description: job.description || '',
        responsibilities: job.responsibilities || '',
        benefits: job.benefits || '',
        is_active: job.is_active,
      });
      setEditingId(id);
      setShowForm(true);
      await rootStore.fetchJobSkills(id);
      const jobSkills = rootStore.jobSkills[id] || [];
      setSelectedSkillIds(jobSkills.map((s) => s.id));
    }
  };

  const handleDelete = async (id: number) => {
    if (confirm('Удалить вакансию?')) {
      await rootStore.deleteJob(id);
    }
  };

  const toggleSkill = (skillId: number) => {
    setSelectedSkillIds((prev) =>
      prev.includes(skillId) ? prev.filter((id) => id !== skillId) : [...prev, skillId]
    );
  };

  const getCompanyName = (companyId: number) => {
    return companies.find((c) => c.id === companyId)?.name || 'Неизвестно';
  };

  // Group skills by category for the checkbox list
  const skillsByCategory = skills.reduce<Record<string, typeof skills>>((acc, skill) => {
    if (!acc[skill.category]) acc[skill.category] = [];
    acc[skill.category].push(skill);
    return acc;
  }, {});

  if (loading.jobs || loading.companies) {
    return <div className="loading-spinner">⏳ Загрузка...</div>;
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
        <h1 style={{ fontSize: '2rem', fontWeight: '700' }}>Вакансии</h1>
        <button className="btn btn-primary" onClick={() => setShowForm(!showForm)}>
          {showForm ? '❌ Закрыть' : '➕ Добавить вакансию'}
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}

      {showForm && (
        <div className="card" style={{ marginBottom: '24px' }}>
          <h2 className="card-title">
            {editingId ? 'Редактировать вакансию' : 'Новая вакансия'}
          </h2>
          <form onSubmit={handleSubmit}>
            <div className="grid grid-2">
              <div className="form-group">
                <label className="form-label">Компания *</label>
                <select
                  className="form-select"
                  value={formData.company_id}
                  onChange={(e) => setFormData({ ...formData, company_id: Number(e.target.value) })}
                  required
                >
                  <option value={0}>Выберите компанию</option>
                  {companies.map((c) => (
                    <option key={c.id} value={c.id}>
                      {c.name}
                    </option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label className="form-label">Название *</label>
                <input
                  className="form-input"
                  type="text"
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                  required
                />
              </div>
              <div className="form-group">
                <label className="form-label">Уровень *</label>
                <select
                  className="form-select"
                  value={formData.level}
                  onChange={(e) => setFormData({ ...formData, level: e.target.value })}
                  required
                >
                  <option value="Junior">Junior</option>
                  <option value="Middle">Middle</option>
                  <option value="Senior">Senior</option>
                  <option value="Lead">Lead</option>
                  <option value="Другое">Другое</option>
                </select>
              </div>
              <div className="form-group">
                <label className="form-label">Специализация</label>
                <select
                  className="form-select"
                  value={formData.specialization}
                  onChange={(e) => setFormData({ ...formData, specialization: e.target.value })}
                >
                  <option value="">— Не указано —</option>
                  {skills.map((s) => (
                    <option key={s.id} value={s.name}>
                      {s.name}
                    </option>
                  ))}
                </select>
              </div>
              <div className="form-group">
                <label className="form-label">Минимальная зарплата</label>
                <input
                  className="form-input"
                  type="number"
                  value={formData.salary_min || ''}
                  onChange={(e) => setFormData({ ...formData, salary_min: Number(e.target.value) || undefined })}
                />
              </div>
              <div className="form-group">
                <label className="form-label">Максимальная зарплата</label>
                <input
                  className="form-input"
                  type="number"
                  value={formData.salary_max || ''}
                  onChange={(e) => setFormData({ ...formData, salary_max: Number(e.target.value) || undefined })}
                />
              </div>
              <div className="form-group">
                <label className="form-label">Валюта</label>
                <select
                  className="form-select"
                  value={formData.salary_currency}
                  onChange={(e) => setFormData({ ...formData, salary_currency: e.target.value })}
                >
                  <option value="RUB">RUB</option>
                  <option value="USD">USD</option>
                  <option value="EUR">EUR</option>
                </select>
              </div>
              <div className="form-group">
                <label className="form-label">Опыт работы</label>
                <input
                  className="form-input"
                  type="text"
                  value={formData.experience_years}
                  onChange={(e) => setFormData({ ...formData, experience_years: e.target.value })}
                  placeholder="например: 3+ лет"
                />
              </div>
              <div className="form-group">
                <label className="form-label">Локация</label>
                <input
                  className="form-input"
                  type="text"
                  value={formData.location}
                  onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                />
              </div>
              <div className="form-group">
                <label className="form-checkbox">
                  <input
                    type="checkbox"
                    checked={formData.remote_available}
                    onChange={(e) => setFormData({ ...formData, remote_available: e.target.checked })}
                  />
                  <span>Удаленная работа</span>
                </label>
              </div>
            </div>

            <div className="form-group">
              <label className="form-label">Описание</label>
              <textarea
                className="form-textarea"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              />
            </div>

            {/* Skills multi-select */}
            <div className="form-group">
              <label className="form-label">
                Навыки{selectedSkillIds.length > 0 && (
                  <span style={{ marginLeft: '8px', color: '#6b7280', fontWeight: '400' }}>
                    выбрано: {selectedSkillIds.length}
                  </span>
                )}
              </label>
              {skills.length === 0 ? (
                <div style={{ color: '#6b7280', fontSize: '14px' }}>Нет доступных навыков</div>
              ) : (
                <div
                  style={{
                    border: '1px solid #e5e7eb',
                    borderRadius: '8px',
                    padding: '12px 16px',
                    maxHeight: '220px',
                    overflowY: 'auto',
                    display: 'flex',
                    flexDirection: 'column',
                    gap: '12px',
                  }}
                >
                  {Object.entries(skillsByCategory).map(([category, catSkills]) => (
                    <div key={category}>
                      <div style={{ fontSize: '12px', fontWeight: '600', color: '#6b7280', marginBottom: '6px', textTransform: 'uppercase' }}>
                        {category}
                      </div>
                      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
                        {catSkills.map((skill) => {
                          const checked = selectedSkillIds.includes(skill.id);
                          return (
                            <label
                              key={skill.id}
                              style={{
                                display: 'flex',
                                alignItems: 'center',
                                gap: '6px',
                                cursor: 'pointer',
                                padding: '4px 10px',
                                borderRadius: '6px',
                                border: `1px solid ${checked ? '#4f46e5' : '#e5e7eb'}`,
                                background: checked ? '#eef2ff' : '#fff',
                                fontSize: '13px',
                                userSelect: 'none',
                              }}
                            >
                              <input
                                type="checkbox"
                                checked={checked}
                                onChange={() => toggleSkill(skill.id)}
                                style={{ margin: 0 }}
                              />
                              {skill.name}
                            </label>
                          );
                        })}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <div style={{ display: 'flex', gap: '12px' }}>
              <button type="submit" className="btn btn-primary">
                {editingId ? '💾 Сохранить' : '➕ Создать'}
              </button>
              <button type="button" className="btn btn-secondary" onClick={resetForm}>
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
              <th>Компания</th>
              <th>Уровень</th>
              <th>Зарплата</th>
              <th>Навыки</th>
              <th>Удаленка</th>
              <th>Статус</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            {jobs.map((job) => {
              const jobSkills = rootStore.jobSkills[job.id] || [];
              return (
                <tr key={job.id}>
                  <td>{job.id}</td>
                  <td style={{ fontWeight: '600' }}>{job.title}</td>
                  <td>{getCompanyName(job.company_id)}</td>
                  <td>
                    <span className="badge badge-primary">{job.level}</span>
                  </td>
                  <td>
                    {job.salary_min && job.salary_max
                      ? `${job.salary_min.toLocaleString()} - ${job.salary_max.toLocaleString()} ${job.salary_currency}`
                      : job.salary_min
                      ? `от ${job.salary_min.toLocaleString()} ${job.salary_currency}`
                      : '—'}
                  </td>
                  <td>
                    {jobSkills.length > 0 ? (
                      <div style={{ display: 'flex', gap: '4px', flexWrap: 'wrap', maxWidth: '220px' }}>
                        {jobSkills.map((s) => (
                          <span key={s.id} className="badge badge-primary" style={{ fontSize: '11px' }}>
                            {s.name}
                          </span>
                        ))}
                      </div>
                    ) : (
                      <span style={{ color: '#9ca3af' }}>—</span>
                    )}
                  </td>
                  <td>{job.remote_available ? '✅' : '❌'}</td>
                  <td>
                    <span className={`badge ${job.is_active ? 'badge-success' : 'badge-warning'}`}>
                      {job.is_active ? 'Активна' : 'Неактивна'}
                    </span>
                  </td>
                  <td>
                    <div style={{ display: 'flex', gap: '8px' }}>
                      <button
                        className="btn btn-secondary"
                        onClick={() => handleEdit(job.id)}
                        style={{ padding: '6px 12px', fontSize: '12px' }}
                      >
                        ✏️
                      </button>
                      <button
                        className="btn btn-danger"
                        onClick={() => handleDelete(job.id)}
                        style={{ padding: '6px 12px', fontSize: '12px' }}
                      >
                        🗑️
                      </button>
                    </div>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
        {jobs.length === 0 && (
          <div style={{ textAlign: 'center', padding: '40px', color: '#6b7280' }}>
            Вакансий пока нет
          </div>
        )}
      </div>
    </div>
  );
});
