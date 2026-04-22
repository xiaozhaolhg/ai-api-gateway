import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/Layout';
import Providers from './pages/Providers';
import Users from './pages/Users';
import APIKeys from './pages/APIKeys';
import Usage from './pages/Usage';
import Health from './pages/Health';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Navigate to="/providers" replace />} />
        <Route path="/providers" element={
          <Layout>
            <Providers />
          </Layout>
        } />
        <Route path="/users" element={
          <Layout>
            <Users />
          </Layout>
        } />
        <Route path="/api-keys" element={
          <Layout>
            <APIKeys />
          </Layout>
        } />
        <Route path="/usage" element={
          <Layout>
            <Usage />
          </Layout>
        } />
        <Route path="/health" element={
          <Layout>
            <Health />
          </Layout>
        } />
      </Routes>
    </Router>
  );
}

export default App;
