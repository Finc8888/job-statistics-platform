import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Layout } from './components/Layout';
import { Dashboard } from './pages/Dashboard';
import { Companies } from './pages/Companies';
import { Jobs } from './pages/Jobs';
import { Skills } from './pages/Skills';
import { Statistics } from './pages/Statistics';
import { styles } from './styles';

export const App = () => {
  return (
    <>
      <style>{styles}</style>
      <Router>
        <Layout>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/companies" element={<Companies />} />
            <Route path="/jobs" element={<Jobs />} />
            <Route path="/skills" element={<Skills />} />
            <Route path="/statistics" element={<Statistics />} />
          </Routes>
        </Layout>
      </Router>
    </>
  );
};
