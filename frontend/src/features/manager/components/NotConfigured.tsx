import React from 'react';
import { useTranslation } from 'react-i18next';

import '@features/manager/components/notconfigured.css';

interface NotConfiguredProps {}

const NotConfigured: React.FC<NotConfiguredProps> = () => {
  const { t } = useTranslation();
  return (
    <div className="notconfigured-container">
      <img
        src="/notconfigured.svg"
        alt="Configuration Error"
        className="notconfigured-container__icon"
      />
      <div className="notconfigured-container__title">
        <span className="notconfigured-container__title-text">
          {t('manager.notconfigured_title', {
            fallback: 'Configuration required',
          })}
        </span>
      </div>
      <div className="notconfigured-container__message">
        {t('manager.notconfigured_message', {
          fallback:
            'Please contact your administrator to configure the document server',
        })}
      </div>
    </div>
  );
};

export default NotConfigured;
