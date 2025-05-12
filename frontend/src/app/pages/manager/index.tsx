import React from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { useFilesStore } from '@features/file/stores/useFileStore';
import { Empty } from '@features/file/components/Empty';
import { FilesList } from '@features/file/components/List';
import { Searchbar } from '@features/file/components/Search';
import { ServerConfigError } from '@features/file/components/ServerConfigError';
import { Installation } from '@features/file/components/Installation';
import { Button } from '@components/Button';
import { Layout } from '@components/Layout';

import "@app/pages/manager/index.css";

export const ManagerPage = () => {
  const { t } = useTranslation();
  const {
    searchQuery,
    filteredDocuments,
    documents,
    loading,
    authError,
    serverConfigError,
  } = useFilesStore();

  const renderContent = () => {
    if (authError) {
      return <Installation />;
    }
    
    if (serverConfigError) {
      return <ServerConfigError />;
    }

    if (loading && documents.length === 0) {
      return <div className="manager-container__main__loading">{t('manager.loading', { fallback: 'Loading...' })}</div>;
    }

    if (documents.length === 0) {
      return <Empty />;
    }
    
    if (searchQuery !== '' && filteredDocuments.length === 0) {
      return <div>{t('manager.notfound', { searchQuery })}</div>;
    }
    
    return <FilesList />;
  };

  return (
    <Layout title={t('manager.title')} subtitle={!authError && !serverConfigError && documents.length > 0 ? t('manager.subtitle') : ''} footerText={t('manager.footer')}>
      <div className='manager-container'>
        <div className='manager-container_shifted'>
          {!authError && !serverConfigError && documents.length > 0 && <Searchbar />}
        </div>
        <div className='manager-container__main'>
          <div className='manager-container_shifted manager-container__main__files'>
            {renderContent()}
          </div>
        </div>
        <div className='manager-container_shifted'>
          {!authError && !serverConfigError && (
            <Link 
              to="/create"
              state={{ isBack: false }}
            >
              <Button name={t('manager.create')} variant='primary' />
            </Link>
          )}
        </div>
      </div>
    </Layout>
  );
};
