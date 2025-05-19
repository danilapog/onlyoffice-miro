import React from 'react';
import { useTranslation } from 'react-i18next';

import Layout from '@components/Layout';

import Creator from '@features/manager/components/Creator';

import '@app/pages/creation/index.css';

const CreationPage: React.FC = () => {
  const { t } = useTranslation();
  return (
    <Layout
      title={t('pages.creation.title')}
      subtitle={t('pages.creation.subtitle')}
      footerText={t('pages.creation.footer')}
      backTo="/"
    >
      <Creator />
    </Layout>
  );
};

export default CreationPage;
