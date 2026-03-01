import { makeAutoObservable, runInAction } from 'mobx';
import {
  companiesApi,
  jobsApi,
  jobSkillsApi,
  skillsApi,
  locationsApi,
  statsApi,
} from '../services/api';
import type {
  Company,
  Job,
  Skill,
  Location,
  TopSkill,
  SkillSalary,
  CompanyStats,
  DatabaseStats,
  LanguageStats,
  CompanyForm,
  JobForm,
  SkillForm,
  LocationForm,
} from '../types';

class RootStore {
  // Data
  companies: Company[] = [];
  jobs: Job[] = [];
  skills: Skill[] = [];
  locations: Location[] = [];
  jobSkills: Record<number, Skill[]> = {};

  // Statistics
  topSkills: TopSkill[] = [];
  skillSalaries: SkillSalary[] = [];
  companyStats: CompanyStats[] = [];
  databaseStats: DatabaseStats[] = [];
  languageStats: LanguageStats[] = [];

  // Loading states
  loading = {
    companies: false,
    jobs: false,
    skills: false,
    locations: false,
    stats: false,
  };

  // Error state
  error: string | null = null;

  constructor() {
    makeAutoObservable(this);
  }

  // Companies
  async fetchCompanies() {
    this.loading.companies = true;
    this.error = null;
    try {
      const response = await companiesApi.getAll();
      runInAction(() => {
        this.companies = response.data;
        this.loading.companies = false;
      });
    } catch (error) {
      runInAction(() => {
        this.error = 'Ошибка загрузки компаний';
        this.loading.companies = false;
      });
    }
  }

  async createCompany(data: CompanyForm) {
    try {
      const response = await companiesApi.create(data);
      runInAction(() => {
        this.companies.push(response.data);
      });
      return response.data;
    } catch (error) {
      this.error = 'Ошибка создания компании';
      throw error;
    }
  }

  async updateCompany(id: number, data: CompanyForm) {
    try {
      const response = await companiesApi.update(id, data);
      runInAction(() => {
        const index = this.companies.findIndex((c) => c.id === id);
        if (index !== -1) {
          this.companies[index] = response.data;
        }
      });
      return response.data;
    } catch (error) {
      this.error = 'Ошибка обновления компании';
      throw error;
    }
  }

  async deleteCompany(id: number) {
    try {
      await companiesApi.delete(id);
      runInAction(() => {
        this.companies = this.companies.filter((c) => c.id !== id);
      });
    } catch (error) {
      this.error = 'Ошибка удаления компании';
      throw error;
    }
  }

  // Jobs
  async fetchJobs() {
    this.loading.jobs = true;
    this.error = null;
    try {
      const response = await jobsApi.getAll();
      runInAction(() => {
        this.jobs = response.data;
        this.loading.jobs = false;
      });
    } catch (error) {
      runInAction(() => {
        this.error = 'Ошибка загрузки вакансий';
        this.loading.jobs = false;
      });
    }
  }

  async createJob(data: JobForm) {
    try {
      const response = await jobsApi.create(data);
      runInAction(() => {
        this.jobs.push(response.data);
      });
      return response.data;
    } catch (error) {
      this.error = 'Ошибка создания вакансии';
      throw error;
    }
  }

  async updateJob(id: number, data: JobForm) {
    try {
      const response = await jobsApi.update(id, data);
      runInAction(() => {
        const index = this.jobs.findIndex((j) => j.id === id);
        if (index !== -1) {
          this.jobs[index] = response.data;
        }
      });
      return response.data;
    } catch (error) {
      this.error = 'Ошибка обновления вакансии';
      throw error;
    }
  }

  async deleteJob(id: number) {
    try {
      await jobsApi.delete(id);
      runInAction(() => {
        this.jobs = this.jobs.filter((j) => j.id !== id);
      });
    } catch (error) {
      this.error = 'Ошибка удаления вакансии';
      throw error;
    }
  }

  // Skills
  async fetchSkills() {
    this.loading.skills = true;
    this.error = null;
    try {
      const response = await skillsApi.getAll();
      runInAction(() => {
        this.skills = response.data;
        this.loading.skills = false;
      });
    } catch (error) {
      runInAction(() => {
        this.error = 'Ошибка загрузки навыков';
        this.loading.skills = false;
      });
    }
  }

  async createSkill(data: SkillForm) {
    try {
      const response = await skillsApi.create(data);
      runInAction(() => {
        this.skills.push(response.data);
      });
      return response.data;
    } catch (error) {
      this.error = 'Ошибка создания навыка';
      throw error;
    }
  }

  async updateSkill(id: number, data: SkillForm) {
    try {
      const response = await skillsApi.update(id, data);
      runInAction(() => {
        const index = this.skills.findIndex((s) => s.id === id);
        if (index !== -1) {
          this.skills[index] = response.data;
        }
      });
      return response.data;
    } catch (error) {
      this.error = 'Ошибка обновления навыка';
      throw error;
    }
  }

  async deleteSkill(id: number) {
    try {
      await skillsApi.delete(id);
      runInAction(() => {
        this.skills = this.skills.filter((s) => s.id !== id);
      });
    } catch (error) {
      this.error = 'Ошибка удаления навыка';
      throw error;
    }
  }

  // Job Skills
  async fetchJobSkills(jobId: number) {
    try {
      const response = await jobSkillsApi.getByJobId(jobId);
      runInAction(() => {
        this.jobSkills[jobId] = response.data;
      });
    } catch {
      // silent fail — skills are non-critical
    }
  }

  async fetchAllJobSkills() {
    if (this.jobs.length === 0) return;
    const results = await Promise.all(
      this.jobs.map((job) =>
        jobSkillsApi.getByJobId(job.id).then((r) => ({ jobId: job.id, skills: r.data }))
      )
    );
    runInAction(() => {
      results.forEach(({ jobId, skills }) => {
        this.jobSkills[jobId] = skills;
      });
    });
  }

  async setJobSkills(jobId: number, skillIds: number[]) {
    try {
      const response = await jobSkillsApi.setForJob(jobId, skillIds);
      runInAction(() => {
        this.jobSkills[jobId] = response.data;
      });
    } catch (error) {
      this.error = 'Ошибка сохранения навыков';
      throw error;
    }
  }

  // Locations
  async fetchLocations() {
    this.loading.locations = true;
    this.error = null;
    try {
      const response = await locationsApi.getAll();
      runInAction(() => {
        this.locations = response.data;
        this.loading.locations = false;
      });
    } catch (error) {
      runInAction(() => {
        this.error = 'Ошибка загрузки локаций';
        this.loading.locations = false;
      });
    }
  }

  async createLocation(data: LocationForm) {
    try {
      const response = await locationsApi.create(data);
      runInAction(() => {
        this.locations.push(response.data);
      });
      return response.data;
    } catch (error) {
      this.error = 'Ошибка создания локации';
      throw error;
    }
  }

  async deleteLocation(id: number) {
    try {
      await locationsApi.delete(id);
      runInAction(() => {
        this.locations = this.locations.filter((l) => l.id !== id);
      });
    } catch (error) {
      this.error = 'Ошибка удаления локации';
      throw error;
    }
  }

  // Statistics
  async fetchAllStats() {
    this.loading.stats = true;
    this.error = null;
    try {
      const [topSkills, skillSalaries, companyStats, databaseStats, languageStats] =
        await Promise.all([
          statsApi.getTopSkills(10),
          statsApi.getSkillSalaries(1),
          statsApi.getCompanyStats(),
          statsApi.getDatabaseStats(),
          statsApi.getLanguageStats(),
        ]);

      runInAction(() => {
        this.topSkills = topSkills.data;
        this.skillSalaries = skillSalaries.data;
        this.companyStats = companyStats.data;
        this.databaseStats = databaseStats.data;
        this.languageStats = languageStats.data;
        this.loading.stats = false;
      });
    } catch (error) {
      runInAction(() => {
        this.error = 'Ошибка загрузки статистики';
        this.loading.stats = false;
      });
    }
  }

  clearError() {
    this.error = null;
  }
}

export const rootStore = new RootStore();
