import React from 'react';
import { useTranslation } from 'react-i18next';

import Button from '@components/Button';

import '@features/manager/components/installation.css';

const Installation = () => {
  const { t } = useTranslation();
  return (
    <div className="installation-container">
      <img
        src="/notconfigured.svg"
        alt="Configuration Error"
        className="installation-container__icon"
      />
      <div className="installation-container__title">
        <span className="installation-container__title-text">
          {t('features.manager.installation.error')}
        </span>
      </div>
      <div className="installation-container__message">
        {t('features.manager.installation.description')}
      </div>
      <Button
        name={t('features.manager.installation.button')}
        variant="primary"
        onClick={() => {
          window.open(import.meta.env.VITE_MIRO_INSTALLATION_URL, '_blank');
          window.miro?.board.ui.closePanel();
        }}
        className="installation-container__button"
        title={t('features.manager.installation.button')}
      />
    </div>
  );
};

export default Installation;
