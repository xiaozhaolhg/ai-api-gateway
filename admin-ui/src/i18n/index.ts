import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import enCommon from './locales/en/common.json';
import zhCommon from './locales/zh/common.json';
import enDashboard from './locales/en/dashboard.json';
import zhDashboard from './locales/zh/dashboard.json';
import enProviders from './locales/en/providers.json';
import zhProviders from './locales/zh/providers.json';
import enRouting from './locales/en/routing.json';
import zhRouting from './locales/zh/routing.json';
import enUsers from './locales/en/users.json';
import zhUsers from './locales/zh/users.json';
import enGroups from './locales/en/groups.json';
import zhGroups from './locales/zh/groups.json';
import enTiers from './locales/en/tiers.json';
import zhTiers from './locales/zh/tiers.json';
import enApiKeys from './locales/en/apiKeys.json';
import zhApiKeys from './locales/zh/apiKeys.json';
import enPermissions from './locales/en/permissions.json';
import zhPermissions from './locales/zh/permissions.json';
import enUsage from './locales/en/usage.json';
import zhUsage from './locales/zh/usage.json';
import enBudgets from './locales/en/budgets.json';
import zhBudgets from './locales/zh/budgets.json';
import enPricing from './locales/en/pricing.json';
import zhPricing from './locales/zh/pricing.json';
import enHealth from './locales/en/health.json';
import zhHealth from './locales/zh/health.json';
import enAlerts from './locales/en/alerts.json';
import zhAlerts from './locales/zh/alerts.json';

const resources = {
  en: {
    common: enCommon,
    dashboard: enDashboard,
    providers: enProviders,
    routing: enRouting,
    users: enUsers,
    groups: enGroups,
    tiers: enTiers,
    apiKeys: enApiKeys,
    permissions: enPermissions,
    usage: enUsage,
    budgets: enBudgets,
    pricing: enPricing,
    health: enHealth,
    alerts: enAlerts,
  },
  zh: {
    common: zhCommon,
    dashboard: zhDashboard,
    providers: zhProviders,
    routing: zhRouting,
    users: zhUsers,
    groups: zhGroups,
    tiers: zhTiers,
    apiKeys: zhApiKeys,
    permissions: zhPermissions,
    usage: zhUsage,
    budgets: zhBudgets,
    pricing: zhPricing,
    health: zhHealth,
    alerts: zhAlerts,
  },
};

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'en',
    lng: localStorage.getItem('i18nextLng') || 'en',
    detection: {
      order: ['localStorage', 'navigator'],
      caches: ['localStorage'],
    },
    interpolation: {
      escapeValue: false,
    },
    ns: ['common', 'dashboard', 'providers', 'routing', 'users', 'groups', 'tiers', 'apiKeys', 'permissions', 'usage', 'budgets', 'pricing', 'health', 'alerts'],
    defaultNS: 'common',
  });

export default i18n;
