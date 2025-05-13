import React from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { useFilesStore } from '@features/file/stores/useFileStore';
import { Empty } from '@features/file/components/Empty';
import { FilesList } from '@features/file/components/List';
import { Searchbar } from '@features/file/components/Search';
import { NotConfigured } from '@features/file/components/NotConfigured';
import { Installation } from '@features/file/components/Installation';
import { Button } from '@components/Button';
import { Layout } from '@components/Layout';
import { Spinner } from '@components/Spinner';

import "@app/pages/manager/index.css";
import { useApplicationStore } from '@stores/useApplicationStore';
import { SettingsPage } from '../settings';

export const ManagerPage = () => {
  const { t } = useTranslation();
  const { admin } = useApplicationStore();
  const {
    searchQuery,
    filteredDocuments,
    documents,
    loading,
    authError,
    serverConfigError,
    initialized,
  } = useFilesStore();

  const isInitialLoading = loading && !initialized;

  if (serverConfigError && admin)
    return <SettingsPage forceDisableBack={true} />;

  const renderContent = () => {
    if (authError) {
      return <Installation />;
    }
    
    if (serverConfigError && !admin) {
      return <NotConfigured />;
    }

    if (loading && documents.length === 0) {
      return <div className="manager-container__main__loading">
        <Spinner size="large" />
      </div>;
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
    <Layout 
      title={serverConfigError ? '' : t('manager.title')} 
      subtitle={!authError && !serverConfigError && documents.length > 0 ? t('manager.subtitle') : ''} 
      footerText={t('manager.footer')}
      settings={!serverConfigError}
    >
      <div className='manager-container'>
        <div className='manager-container_shifted'>
          {!authError && !serverConfigError && documents.length > 0 && <Searchbar />}
        </div>
        <div className='manager-container__main'>
          <div className='manager-container_shifted manager-container__main__files'>
            {renderContent()}
          </div>
        </div>
        <div className='manager-container_shifted manager-container__button-container'>
          {!authError && !serverConfigError && (
            <Link 
              to="/create"
              state={{ isBack: false }}
              className={isInitialLoading ? 'disabled-link' : ''}
              onClick={e => isInitialLoading && e.preventDefault()}
            >
              <Button 
                name={t('manager.create')} 
                variant='primary' 
                disabled={isInitialLoading}
              />
            </Link>
          )}
        </div>
      </div>
    </Layout>
  );
};
