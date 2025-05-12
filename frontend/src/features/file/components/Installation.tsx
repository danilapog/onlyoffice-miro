import React from 'react';
import { useTranslation } from 'react-i18next';
import { Button } from '@components/Button';

export const Installation = () => {
  const { t } = useTranslation();
  
  return (
    <div className="installation-error">
      <div>{t('manager.installation_error', { fallback: 'App Installation Required' })}</div>
      <div>{t('manager.installation_error_description', { fallback: 'Please install or reinstall the app to continue' })}</div>
      <Button 
        name={t('manager.installation_error_button', { fallback: 'Install' })} 
        variant='primary' 
        onClick={() => {
          window.open(import.meta.env.VITE_MIRO_INSTALLATION_URL, '_blank');
        }} 
      />
    </div>
  );
}; 