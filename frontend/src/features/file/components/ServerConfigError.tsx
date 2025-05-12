import React from 'react';
import { useTranslation } from 'react-i18next';
import { Button } from '@components/Button';
import { useApplicationStore } from '@stores/useApplicationStore';

export const ServerConfigError = () => {
  const { t } = useTranslation();
  const { admin } = useApplicationStore();
  
  return (
    <div className="server-config-error">
      <div>{t('manager.server_config_error', { fallback: 'Document Server Configuration Required' })}</div>
      <div>{t('manager.server_config_error_description', { fallback: 'Please configure your document server in the app settings' })}</div>
      
      {admin ? (
        <Button 
          name={t('manager.server_config_error_button', { fallback: 'Configure' })} 
          variant='primary' 
          onClick={() => {
            window.location.href = '#/settings';
          }} 
        />
      ) : (
        <div className="server-config-error__non-admin">
          {t('manager.server_config_error_non_admin', { 
            fallback: 'Please contact your administrator to configure the document server' 
          })}
        </div>
      )}
    </div>
  );
}; 