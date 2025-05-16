import React from 'react';
import { useTranslation } from 'react-i18next';

import Layout from '@components/Layout';

import Form from '@features/settings/components/Form';

import useSettingsStore from '@features/settings/stores/useSettingsStore';

import '@app/pages/settings/index.css';

interface SettingsPageProps {
  forceDisableBack?: boolean;
}

const SettingsPage: React.FC<SettingsPageProps> = ({
  forceDisableBack = false,
}) => {
  const { t } = useTranslation();
  const { hasSettings } = useSettingsStore();

  const backTo = !forceDisableBack && hasSettings ? '/' : undefined;

  return (
    <Layout
      title={t('settings.title')}
      footerText={t('settings.footer')}
      settings={false}
      help
      backTo={backTo}
    >
      <Form />
    </Layout>
  );
};

SettingsPage.displayName = 'SettingsPage';

export default SettingsPage;
