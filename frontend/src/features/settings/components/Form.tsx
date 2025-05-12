import React, { forwardRef, FormEvent } from 'react';
import { useTranslation } from 'react-i18next';

import { useSettingsStore } from '@features/settings/stores/useSettingsStore';
import { Banner } from '@features/settings/components/Banner';
import { FormInput } from '@components/Input';
import { Button } from '@components/Button';

import '@features/settings/components/form.css';

interface FormProps extends React.HTMLAttributes<HTMLDivElement> { }

export const Form = forwardRef<HTMLDivElement, FormProps>(({
  className,
  children,
  ...props
}, ref) => {
  const { t } = useTranslation();
  const {
    address,
    header,
    secret,
    loading,
    demo,
    demoStarted,
    setAddress,
    setHeader,
    setSecret,
    setDemo,
    saveSettings,
  } = useSettingsStore();

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    await saveSettings();
  };

  const isDemoExpired = demoStarted ? 
    (() => {
      const startTime = new Date(demoStarted).getTime();
      const expiryDays = parseInt(import.meta.env.VITE_ASC_DEMO_EXPIRATION_DAYS || '30', 10);
      const expiryTime = startTime + (expiryDays * 24 * 60 * 60 * 1000);
      const currentTime = Date.now();      
      return currentTime > expiryTime;
    })() : false;

  return (
    <div
      ref={ref}
      className={`form ${className || ''}`}
      {...props}
    >
      <div className="form__content">
        <p className="form__description">
          {t('settings.description')}
        </p>
        <form onSubmit={handleSubmit} className="form__fields" autoComplete="off">
          <div className="form__field">
            <FormInput
              label={t('settings.address')}
              name="address"
              type="text"
              value={address}
              disabled={loading}
              onChange={(e) => setAddress(e.target.value)}
              autoComplete="off"
            />
          </div>
          <div className="form__field">
            <FormInput
              label={t('settings.secret')}
              name="secret"
              type="password"
              value={secret}
              disabled={loading}
              onChange={(e) => setSecret(e.target.value)}
              autoComplete="off"
            />
          </div>
          <div className="form__field">
            <FormInput
              label={t('settings.header')}
              name="header"
              type="text"
              value={header}
              disabled={loading}
              onChange={(e) => setHeader(e.target.value)}
              autoComplete="off"
            />
          </div>
          
          {isDemoExpired && <Banner />}
          
          {!isDemoExpired && (
            <div className="form__checkbox-container">
              <label className="form__checkbox-label">
                <input
                type="checkbox"
                className="form__checkbox"
                checked={demo}
                disabled={loading || !!demoStarted}
                onChange={() => setDemo(!demo)}
              />
              <span className="form__checkbox-text">{t('settings.demo.title')}</span>
            </label>
            <p className="form__checkbox-description">
                {!demoStarted && t('settings.demo.description')}
                {demoStarted && t('settings.demo.started', { date: new Date(demoStarted).toLocaleDateString('en-GB').split('/').join('.') })}
              </p>
            </div>
          )}

          <div className="form__button-container">
            <Button
              type="submit"
              name={t('settings.save')}
              variant="primary"
              disabled={loading}
            />
          </div>
        </form>
      </div>
    </div>
  );
});
