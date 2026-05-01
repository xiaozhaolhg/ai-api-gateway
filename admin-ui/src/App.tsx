import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { I18nextProvider, useTranslation } from 'react-i18next';
import { ConfigProvider, App as AntdApp } from 'antd';
import enUS from 'antd/locale/en_US';
import zhCN from 'antd/locale/zh_CN';
import { AuthProvider } from './contexts/AuthContext';
import { ProtectedRoute } from './components/ProtectedRoute';
import { AppShell } from './components/AppShell';
import DevTools from './components/DevTools';
import { Login } from './pages/Login';
import { Register } from './pages/Register';
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
              <Route path="/register" element={<Register />} />
              <Route
                path="/"
                element={
                  <ProtectedRoute>
                    <AppShell />
                  </ProtectedRoute>
                }
              >
                <Route index element={<Dashboard />} />
                <Route path="providers" element={<ProtectedRoute requiredRole="admin"><Providers /></ProtectedRoute>} />
                <Route path="routing" element={<ProtectedRoute requiredRole="admin"><RoutingRules /></ProtectedRoute>} />
                <Route path="users" element={<ProtectedRoute requiredRole="admin"><Users /></ProtectedRoute>} />
                <Route path="groups" element={<ProtectedRoute requiredRole="admin"><Groups /></ProtectedRoute>} />
                <Route path="api-keys" element={<ProtectedRoute requiredRole="user"><APIKeys /></ProtectedRoute>} />
                <Route path="permissions" element={<ProtectedRoute requiredRole="admin"><Permissions /></ProtectedRoute>} />
                <Route path="usage" element={<Usage />} />
                <Route path="budgets" element={<ProtectedRoute requiredRole="admin"><Budgets /></ProtectedRoute>} />
                <Route path="pricing" element={<ProtectedRoute requiredRole="admin"><PricingRules /></ProtectedRoute>} />
                <Route path="health" element={<Health />} />
                <Route path="alerts" element={<ProtectedRoute requiredRole="admin"><Alerts /></ProtectedRoute>} />
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
      <DevTools />
    </I18nextProvider>
  );
}

export default App;
