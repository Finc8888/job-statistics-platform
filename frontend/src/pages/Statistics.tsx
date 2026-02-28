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
  LineChart,
  Line,
  RadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  Radar,
} from 'recharts';

export const Statistics = observer(() => {
  useEffect(() => {
    rootStore.fetchAllStats();
  }, []);

  const {
    topSkills,
    skillSalaries,
    companyStats,
    databaseStats,
    languageStats,
    loading,
  } = rootStore;

  if (loading.stats) {
    return <div className="loading-spinner">⏳ Загрузка статистики...</div>;
  }

  // Database stats data
  const dbChartData = databaseStats.slice(0, 8).map((db) => ({
    name: db.database,
    'Вакансий': db.vacancies,
    'Обязательно': db.required,
    'Желательно': db.nice_to_have,
  }));

  // Skill salaries data
  const salaryChartData = skillSalaries.slice(0, 10).map((s) => ({
    name: s.name,
    'Средняя зарплата (₽)': Math.round(s.avg_salary),
    'Вакансий': s.vacancy_count,
  }));

  // Company stats data
  const companyChartData = companyStats.map((c) => ({
    name: c.company,
    'Вакансий': c.vacancies_count,
    'Мин. зарплата': Math.round(c.min_salary / 1000),
    'Макс. зарплата': Math.round(c.max_salary / 1000),
  }));

  // Top skills radar data
  const radarData = topSkills.slice(0, 6).map((skill) => ({
    skill: skill.name,
    value: skill.vacancy_count,
  }));

  return (
    <div>
      <h1 style={{ fontSize: '2rem', fontWeight: '700', marginBottom: '32px' }}>
        📊 Детальная статистика
      </h1>

      {/* Top Skills */}
      <div className="card" style={{ marginBottom: '24px' }}>
        <h2 className="card-title">Топ-10 навыков по востребованности</h2>
        <ResponsiveContainer width="100%" height={400}>
          <BarChart data={topSkills.slice(0, 10).map((s) => ({
            name: s.name,
            'Всего': s.vacancy_count,
            'Обязательно': s.required_count,
            'Желательно': s.nice_to_have_count,
          }))}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="name" angle={-45} textAnchor="end" height={100} />
            <YAxis />
            <Tooltip />
            <Legend />
            <Bar dataKey="Всего" fill="#3b82f6" />
            <Bar dataKey="Обязательно" fill="#10b981" />
            <Bar dataKey="Желательно" fill="#f59e0b" />
          </BarChart>
        </ResponsiveContainer>
      </div>

      {/* Salaries */}
      <div className="card" style={{ marginBottom: '24px' }}>
        <h2 className="card-title">Средняя зарплата по навыкам (топ-10)</h2>
        <ResponsiveContainer width="100%" height={400}>
          <LineChart data={salaryChartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="name" angle={-45} textAnchor="end" height={100} />
            <YAxis yAxisId="left" />
            <YAxis yAxisId="right" orientation="right" />
            <Tooltip formatter={(value: number) => value.toLocaleString('ru-RU')} />
            <Legend />
            <Line
              yAxisId="left"
              type="monotone"
              dataKey="Средняя зарплата (₽)"
              stroke="#10b981"
              strokeWidth={2}
              dot={{ r: 5 }}
            />
            <Line
              yAxisId="right"
              type="monotone"
              dataKey="Вакансий"
              stroke="#3b82f6"
              strokeWidth={2}
              dot={{ r: 5 }}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>

      <div className="grid grid-2">
        {/* Database Stats */}
        <div className="card">
          <h2 className="card-title">Востребованность баз данных</h2>
          <ResponsiveContainer width="100%" height={350}>
            <BarChart data={dbChartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" angle={-45} textAnchor="end" height={80} />
              <YAxis />
              <Tooltip />
              <Legend />
              <Bar dataKey="Вакансий" fill="#3b82f6" />
              <Bar dataKey="Обязательно" fill="#10b981" />
            </BarChart>
          </ResponsiveContainer>
        </div>

        {/* Skills Radar */}
        <div className="card">
          <h2 className="card-title">Топ-6 навыков (радар)</h2>
          <ResponsiveContainer width="100%" height={350}>
            <RadarChart data={radarData}>
              <PolarGrid />
              <PolarAngleAxis dataKey="skill" />
              <PolarRadiusAxis />
              <Radar
                name="Вакансий"
                dataKey="value"
                stroke="#3b82f6"
                fill="#3b82f6"
                fillOpacity={0.6}
              />
              <Tooltip />
            </RadarChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Company Stats */}
      <div className="card" style={{ marginTop: '24px' }}>
        <h2 className="card-title">Статистика по компаниям</h2>
        <ResponsiveContainer width="100%" height={350}>
          <BarChart data={companyChartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="name" />
            <YAxis yAxisId="left" />
            <YAxis yAxisId="right" orientation="right" />
            <Tooltip />
            <Legend />
            <Bar yAxisId="left" dataKey="Вакансий" fill="#3b82f6" />
            <Bar yAxisId="right" dataKey="Мин. зарплата" fill="#10b981" />
            <Bar yAxisId="right" dataKey="Макс. зарплата" fill="#f59e0b" />
          </BarChart>
        </ResponsiveContainer>
        <p style={{ fontSize: '12px', color: '#6b7280', marginTop: '8px' }}>
          * Зарплаты указаны в тысячах рублей
        </p>
      </div>

      {/* Language Stats Table */}
      <div className="card" style={{ marginTop: '24px' }}>
        <h2 className="card-title">Языки программирования: детальная статистика</h2>
        <table className="table">
          <thead>
            <tr>
              <th>Язык</th>
              <th>Вакансий</th>
              <th>Средняя зарплата</th>
              <th>Компании</th>
              <th>Уровни</th>
            </tr>
          </thead>
          <tbody>
            {languageStats.map((lang) => (
              <tr key={lang.language}>
                <td style={{ fontWeight: '600' }}>{lang.language}</td>
                <td>
                  <span className="badge badge-primary">{lang.vacancy_count}</span>
                </td>
                <td>{Math.round(lang.avg_salary_rub).toLocaleString('ru-RU')} ₽</td>
                <td style={{ fontSize: '12px', maxWidth: '200px' }}>
                  {lang.companies}
                </td>
                <td style={{ fontSize: '12px' }}>{lang.levels}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
});
