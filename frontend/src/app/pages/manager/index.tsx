import React from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { Button } from '@components/Button';
import { Layout } from '@components/Layout';
import { Spinner } from '@components/Spinner';

import { Empty } from '@features/manager/components/Empty';
import { FilesList } from '@features/file/components/List';
import { Installation } from '@features/manager/components/Installation';
import { NotConfigured } from '@features/manager/components/NotConfigured';
import { Searchbar } from '@features/file/components/Search';

import { useFilesStore } from '@features/file/stores/useFileStore';
import { useApplicationStore } from '@stores/useApplicationStore';

import { SettingsPage } from '@app/pages/settings';

import "@app/pages/manager/index.css";

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
    if (authError)
      return <Installation />;
    
    if (serverConfigError && !admin)
      return <NotConfigured />;

    if (loading && documents.length === 0)
      return <div className="manager-container__main__loading">
        <Spinner size="large" />
      </div>;


    if (documents.length === 0)
      return <Empty 
        title={t('empty.title')}
        subtitle={t('empty.subtitle')}
      />;

    if (searchQuery !== '' && filteredDocuments.length === 0)
      return <Empty 
        title={t('empty.search_title')}
        subtitle={t('empty.search_subtitle')}
      />;
    
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
