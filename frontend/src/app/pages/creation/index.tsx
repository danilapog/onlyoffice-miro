import React from 'react';
import { useTranslation } from 'react-i18next';

import { Layout } from '@components/Layout';

import { Creator } from '@features/manager/components/Creator';

import '@app/pages/creation/index.css';

export const CreationPage: React.FC = () => {
  const { t } = useTranslation();
  return (
    <Layout 
      title={t('creation.title')}
      subtitle={t('creation.subtitle')}
      footerText={t('creation.footer')}
      backTo='/'
    >
      <Creator />
    </Layout>
  );
};

