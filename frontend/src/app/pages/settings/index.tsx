import React from 'react';
import { useTranslation } from 'react-i18next';
import { useSettingsStore } from '@features/settings/stores/useSettingsStore';
import { Form } from '@features/settings/components/Form';
import { Layout } from '@components/Layout';

import '@app/pages/settings/index.css';

interface SettingsPageProps {
  forceDisableBack?: boolean;
}

export const SettingsPage: React.FC<SettingsPageProps> = ({ forceDisableBack = false }) => {
  const { t } = useTranslation();
  const { hasSettings } = useSettingsStore();
  
  const backTo = !forceDisableBack && hasSettings ? '/' : undefined;

  return (
    <Layout title={t('settings.title')} subtitle={t('settings.subtitle')} footerText={t('settings.footer')} settings={false} help={true} backTo={backTo}>
      <Form />
    </Layout>
  );
};
