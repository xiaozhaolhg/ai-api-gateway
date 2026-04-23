import React from 'react';
import { Dropdown, Button } from 'antd';
import { GlobalOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';

export const LanguageSwitcher: React.FC = () => {
  const { i18n } = useTranslation();

  const handleLanguageChange = (lang: string) => {
    i18n.changeLanguage(lang);
  };

  const items = [
    {
      key: 'en',
      label: 'English',
      onClick: () => handleLanguageChange('en'),
    },
    {
      key: 'zh',
      label: '中文',
      onClick: () => handleLanguageChange('zh'),
    },
  ];

  return (
    <Dropdown menu={{ items }} placement="bottomRight">
      <Button icon={<GlobalOutlined />} type="text">
        {i18n.language === 'zh' ? '中文' : 'English'}
      </Button>
    </Dropdown>
  );
};
