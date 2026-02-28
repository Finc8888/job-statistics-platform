import { useEffect } from 'react';
import { observer } from 'mobx-react-lite';
import { rootStore } from '../stores/RootStore';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts';

const COLORS = ['#3b82f6', '#8b5cf6', '#10b981', '#f59e0b', '#ef4444', '#06b6d4'];

export const Dashboard = observer(() => {
  useEffect(() => {
    rootStore.fetchCompanies();
    rootStore.fetchJobs();
    rootStore.fetchSkills();
    rootStore.fetchAllStats();
  }, []);

  const { companies, jobs, skills, topSkills, languageStats, loading } = rootStore;

  if (loading.stats || loading.companies || loading.jobs || loading.skills) {
    return <div className="loading-spinner">⏳ Загрузка данных...</div>;
  }

  // Данные для круговой диаграммы уровней
  const levelCounts = jobs.reduce((acc, job) => {
    acc[job.level] = (acc[job.level] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  const levelData = Object.entries(levelCounts).map(([name, value]) => ({
    name,
    value,
  }));

  // Данные для графика топ-10 навыков
  const topSkillsData = topSkills.slice(0, 10).map((skill) => ({
    name: skill.name,
    'Всего вакансий': skill.vacancy_count,
    Обязательно: skill.required_count,
    Желательно: skill.nice_to_have_count,
  }));

  // Данные для графика языков программирования
  const languageChartData = languageStats.slice(0, 8).map((lang) => ({
    name: lang.language,
    'Вакансий': lang.vacancy_count,
    'Ср. зарплата (₽)': Math.round(lang.avg_salary_rub / 1000),
  }));

  return (
    <div>
      <h1 style={{ fontSize: '2rem', fontWeight: '700', marginBottom: '24px' }}>
        Дашборд
      </h1>

      {/* Stats Cards */}
      <div className="grid grid-2" style={{ marginBottom: '32px' }}>
        <div className="stat-card">
          <div className="stat-value" style={{ color: '#3b82f6' }}>
            {companies.length}
          </div>
          <div className="stat-label">Компаний</div>
        </div>
        <div className="stat-card">
          <div className="stat-value" style={{ color: '#8b5cf6' }}>
            {jobs.length}
          </div>
          <div className="stat-label">Вакансий</div>
        </div>
        <div className="stat-card">
          <div className="stat-value" style={{ color: '#10b981' }}>
            {skills.length}
          </div>
          <div className="stat-label">Навыков</div>
        </div>
        <div className="stat-card">
          <div className="stat-value" style={{ color: '#f59e0b' }}>
            {jobs.filter((j) => j.remote_available).length}
          </div>
          <div className="stat-label">Удалённых вакансий</div>
        </div>
      </div>

      {/* Charts */}
      <div className="grid grid-2">
        {/* Top Skills Chart */}
        <div className="card">
          <h2 className="card-title">Топ-10 навыков</h2>
          <ResponsiveContainer width="100%" height={400}>
            <BarChart data={topSkillsData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis
                dataKey="name"
                angle={-45}
                textAnchor="end"
                height={100}
                style={{ fontSize: '12px' }}
              />
              <YAxis />
              <Tooltip />
              <Legend />
              <Bar dataKey="Всего вакансий" fill="#3b82f6" />
              <Bar dataKey="Обязательно" fill="#10b981" />
              <Bar dataKey="Желательно" fill="#f59e0b" />
            </BarChart>
          </ResponsiveContainer>
        </div>

        {/* Job Levels Pie Chart */}
        <div className="card">
          <h2 className="card-title">Распределение по уровням</h2>
          <ResponsiveContainer width="100%" height={400}>
            <PieChart>
              <Pie
                data={levelData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percent }) =>
                  `${name}: ${(percent * 100).toFixed(0)}%`
                }
                outerRadius={120}
                fill="#8884d8"
                dataKey="value"
              >
                {levelData.map((_, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </div>

        {/* Programming Languages Chart */}
        <div className="card" style={{ gridColumn: '1 / -1' }}>
          <h2 className="card-title">Языки программирования: вакансии и зарплаты</h2>
          <ResponsiveContainer width="100%" height={400}>
            <BarChart data={languageChartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" />
              <YAxis yAxisId="left" />
              <YAxis yAxisId="right" orientation="right" />
              <Tooltip />
              <Legend />
              <Bar yAxisId="left" dataKey="Вакансий" fill="#3b82f6" />
              <Bar yAxisId="right" dataKey="Ср. зарплата (₽)" fill="#10b981" />
            </BarChart>
          </ResponsiveContainer>
          <p style={{ fontSize: '12px', color: '#6b7280', marginTop: '8px' }}>
            * Средняя зарплата указана в тысячах рублей
          </p>
        </div>
      </div>
    </div>
  );
});
