import axios from 'axios';
import type {
  Company,
  Job,
  Skill,
  Location,
  TopSkill,
  SkillSalary,
  SkillByLevel,
  CompanyStats,
  DatabaseStats,
  LanguageStats,
  CompanyForm,
  JobForm,
  SkillForm,
  LocationForm,
} from '../types';

const API_BASE_URL = 'http://localhost:8081/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Companies API
export const companiesApi = {
  getAll: () => api.get<Company[]>('/companies'),
  getById: (id: number) => api.get<Company>(`/companies/${id}`),
  create: (data: CompanyForm) => api.post<Company>('/companies', data),
  update: (id: number, data: CompanyForm) => api.put<Company>(`/companies/${id}`, data),
  delete: (id: number) => api.delete(`/companies/${id}`),
};

// Jobs API
export const jobsApi = {
  getAll: () => api.get<Job[]>('/jobs'),
  getById: (id: number) => api.get<Job>(`/jobs/${id}`),
  create: (data: JobForm) => api.post<Job>('/jobs', data),
  update: (id: number, data: JobForm) => api.put<Job>(`/jobs/${id}`, data),
  delete: (id: number) => api.delete(`/jobs/${id}`),
};

// Skills API
export const skillsApi = {
  getAll: () => api.get<Skill[]>('/skills'),
  getById: (id: number) => api.get<Skill>(`/skills/${id}`),
  create: (data: SkillForm) => api.post<Skill>('/skills', data),
  update: (id: number, data: SkillForm) => api.put<Skill>(`/skills/${id}`, data),
  delete: (id: number) => api.delete(`/skills/${id}`),
};

// Locations API
export const locationsApi = {
  getAll: () => api.get<Location[]>('/locations'),
  getByJobId: (jobId: number) => api.get<Location[]>(`/locations/job/${jobId}`),
  create: (data: LocationForm) => api.post<Location>('/locations', data),
  update: (id: number, data: LocationForm) => api.put<Location>(`/locations/${id}`, data),
  delete: (id: number) => api.delete(`/locations/${id}`),
};

// Statistics API
export const statsApi = {
  getTopSkills: (limit = 10) => api.get<TopSkill[]>(`/stats/top-skills?limit=${limit}`),
  getSkillSalaries: (minVacancies = 1) => api.get<SkillSalary[]>(`/stats/skill-salaries?min_vacancies=${minVacancies}`),
  getSkillsByLevel: () => api.get<SkillByLevel[]>('/stats/skills-by-level'),
  getCompanyStats: () => api.get<CompanyStats[]>('/stats/companies'),
  getDatabaseStats: () => api.get<DatabaseStats[]>('/stats/databases'),
  getLanguageStats: () => api.get<LanguageStats[]>('/stats/languages'),
};

export default api;
