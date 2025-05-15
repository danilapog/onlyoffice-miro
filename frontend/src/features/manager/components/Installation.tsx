import React from 'react';
import { useTranslation } from 'react-i18next';

import { Button } from '@components/Button';

import '@features/manager/components/installation.css';

export const Installation = () => {
  const { t } = useTranslation();
  return (
    <div className="installation-container">
      <img src="/notconfigured.svg" alt="Configuration Error" className="installation-container__icon" />
      <div className="installation-container__title">
        <span className="installation-container__title-text">
          {t('manager.installation_error', { 
            fallback: 'App Installation Required' 
          })}
        </span>
      </div>
      <div className="installation-container__message">
        {t('manager.installation_error_description', { 
          fallback: 'Please install or reinstall the app to continue' 
        })}
      </div>
      <Button 
        name={t('manager.installation_error_button', { fallback: 'Install' })} 
        variant='primary' 
        onClick={() => {
          window.open(import.meta.env.VITE_MIRO_INSTALLATION_URL, '_blank');
          window.miro?.board.ui.closePanel();
        }}
        className="installation-container__button"
        title={t('manager.installation_error_button', { fallback: 'Install' })}
      />
    </div>
  );
};