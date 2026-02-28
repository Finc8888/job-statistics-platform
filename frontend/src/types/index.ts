// API Types

export interface Company {
  id: number;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface Job {
  id: number;
  company_id: number;
  title: string;
  level: 'Junior' | 'Middle' | 'Senior' | 'Lead' | 'Другое';
  specialization: string | null;
  salary_min: number | null;
  salary_max: number | null;
  salary_currency: string;
  experience_years: string | null;
  location: string | null;
  remote_available: boolean;
  description: string | null;
  responsibilities: string | null;
  benefits: string | null;
  posted_date: string | null;
  is_active: boolean;
  source_url: string | null;
  created_at: string;
  updated_at: string;
}

export interface Skill {
  id: number;
  name: string;
  category: 'Язык программирования' | 'База данных' | 'Фреймворк' | 'Инструмент' | 'Другое';
  created_at: string;
}

export interface Location {
  id: number;
  job_id: number;
  city: string;
  metro_station: string | null;
  is_primary: boolean;
}

// Statistics Types

export interface TopSkill {
  name: string;
  category: string;
  vacancy_count: number;
  required_count: number;
  nice_to_have_count: number;
}

export interface SkillSalary {
  name: string;
  avg_salary: number;
  vacancy_count: number;
}

export interface SkillByLevel {
  level: string;
  skill_name: string;
  count: number;
}

export interface CompanyStats {
  company: string;
  vacancies_count: number;
  levels: string;
  min_salary: number;
  max_salary: number;
  locations_count: number;
  remote_vacancies: number;
}

export interface DatabaseStats {
  database: string;
  vacancies: number;
  required: number;
  nice_to_have: number;
  companies: string;
  avg_salary: number;
}

export interface LanguageStats {
  language: string;
  vacancy_count: number;
  companies: string;
  levels: string;
  avg_salary_rub: number;
}

// Form Types

export interface CompanyForm {
  name: string;
  description: string;
}

export interface JobForm {
  company_id: number;
  title: string;
  level: string;
  specialization?: string;
  salary_min?: number;
  salary_max?: number;
  salary_currency: string;
  experience_years?: string;
  location?: string;
  remote_available: boolean;
  description?: string;
  responsibilities?: string;
  benefits?: string;
  is_active: boolean;
}

export interface SkillForm {
  name: string;
  category: string;
}

export interface LocationForm {
  job_id: number;
  city: string;
  metro_station?: string;
  is_primary: boolean;
}
