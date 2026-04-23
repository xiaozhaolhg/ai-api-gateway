import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { I18nextProvider, useTranslation } from 'react-i18next';
import { ConfigProvider, App as AntdApp } from 'antd';
import enUS from 'antd/locale/en_US';
import zhCN from 'antd/locale/zh_CN';
import { AuthProvider } from './contexts/AuthContext';
import { ProtectedRoute } from './components/ProtectedRoute';
import { AppShell } from './components/AppShell';
import { Login } from './pages/Login';
import { Dashboard } from './pages/Dashboard';
import Providers from './pages/Providers';
import { RoutingRules } from './pages/RoutingRules';
import Users from './pages/Users';
import { Groups } from './pages/Groups';
import APIKeys from './pages/APIKeys';
import { Permissions } from './pages/Permissions';
import Usage from './pages/Usage';
import { Budgets } from './pages/Budgets';
import { PricingRules } from './pages/PricingRules';
import Health from './pages/Health';
import { Alerts } from './pages/Alerts';
import i18n from './i18n';

function AppContent() {
  const { i18n: i18nHook } = useTranslation();
  const antdLocale = i18nHook.language === 'zh' ? zhCN : enUS;

  return (
    <ConfigProvider locale={antdLocale} key={i18nHook.language}>
      <AntdApp>
        <AuthProvider>
          <Router>
            <Routes>
              <Route path="/login" element={<Login />} />
              <Route
                path="/"
                element={
                  <ProtectedRoute>
                    <AppShell />
                  </ProtectedRoute>
                }
              >
                <Route index element={<Dashboard />} />
                <Route path="providers" element={<Providers />} />
                <Route path="routing" element={<RoutingRules />} />
                <Route path="users" element={<Users />} />
                <Route path="groups" element={<Groups />} />
                <Route path="api-keys" element={<APIKeys />} />
                <Route path="permissions" element={<Permissions />} />
                <Route path="usage" element={<Usage />} />
                <Route path="budgets" element={<Budgets />} />
                <Route path="pricing" element={<PricingRules />} />
                <Route path="health" element={<Health />} />
                <Route path="alerts" element={<Alerts />} />
              </Route>
            </Routes>
          </Router>
        </AuthProvider>
      </AntdApp>
    </ConfigProvider>
  );
}

function App() {
  return (
    <I18nextProvider i18n={i18n}>
      <AppContent />
    </I18nextProvider>
  );
}

export default App;
