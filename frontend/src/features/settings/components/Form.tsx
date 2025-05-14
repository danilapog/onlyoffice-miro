import React, { forwardRef, FormEvent, useState, useEffect, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { validateAddress, validateShortText } from '@utils/validator';

import { Button } from '@components/Button';
import { FormInput } from '@components/Input';

import { Banner } from '@features/settings/components/Banner';

import { useSettingsStore } from '@features/settings/stores/useSettingsStore';
import { useApplicationStore } from '@stores/useApplicationStore';
import { useEmitterStore } from '@stores/useEmitterStore';

import '@features/settings/components/form.css';

interface FormProps extends React.HTMLAttributes<HTMLDivElement> { }

export const Form = forwardRef<HTMLDivElement, FormProps>(({
  className,
  children,
  ...props
}, ref) => {
  const [addressError, setAddressError] = useState('');
  const [secretError, setSecretError] = useState('');
  const [headerError, setHeaderError] = useState('');
  const { t } = useTranslation();
  const navigate = useNavigate();
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
  const { refreshAuthorization } = useApplicationStore();
  const { emitRefreshDocuments } = useEmitterStore();

  const isDemoExpired = demoStarted ? 
    (() => {
      const startTime = new Date(demoStarted).getTime();
      const expiryDays = parseInt(import.meta.env.VITE_ASC_DEMO_EXPIRATION_DAYS || '30', 10);
      const expiryTime = startTime + (expiryDays * 24 * 60 * 60 * 1000);
      const currentTime = Date.now();      
      return currentTime > expiryTime;
    })() : false;

  const hasInputs = address.trim() !== '' || header.trim() !== '' || secret.trim() !== '';
  const fieldsRequired = !demo || isDemoExpired || hasInputs;
  
  const validateAddressField = (value: string): string => {
    if (!fieldsRequired) return '';
    return validateAddress(value) ? '' : t('settings.errors.addressRequired');
  };

  const validateHeaderField = (value: string): string => {
    if (!fieldsRequired) return '';
    return validateShortText(value) ? '' : t('settings.errors.headerRequired');
  };

  const validateSecretField = (value: string): string => {
    if (!fieldsRequired) return '';
    return validateShortText(value) ? '' : t('settings.errors.secretRequired');
  };
  
  const addressErr = useMemo(() => validateAddressField(address), [address, fieldsRequired]);
  const headerErr = useMemo(() => validateHeaderField(header), [header, fieldsRequired]);
  const secretErr = useMemo(() => validateSecretField(secret), [secret, fieldsRequired]);
  
  const hasValidationErrors = !!(addressErr || headerErr || secretErr);
  const saveDisabled = loading || 
    (!hasInputs && !demo) || 
    (!hasInputs && demo && !!demoStarted && !isDemoExpired) || 
    hasValidationErrors;

  useEffect(() => {
    if (fieldsRequired) {
      setAddressError(addressErr);
      setHeaderError(headerErr);
      setSecretError(secretErr);
    } else {
      setAddressError('');
      setHeaderError('');
      setSecretError('');
    }
  }, [demo, isDemoExpired, address, header, secret, saveDisabled, addressErr, headerErr, secretErr]);

  const validateForm = (): boolean => {
    if (!fieldsRequired) return true;
    
    setAddressError(addressErr);
    setHeaderError(headerErr);
    setSecretError(secretErr);
    
    return !addressErr && !headerErr && !secretErr;
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    
    if (validateForm()) {
      await saveSettings();
      await emitRefreshDocuments();
      await refreshAuthorization();
      navigate('/');
    }
  };

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
              error={addressError}
              disabled={loading}
              onChange={(e) => {
                const value = e.target.value;
                setAddress(value);
                setAddressError(validateAddressField(value));
              }}
              required={fieldsRequired}
              autoComplete="off"
            />
          </div>
          <div className="form__field">
            <FormInput
              label={t('settings.secret')}
              name="secret"
              type="password"
              value={secret}
              error={secretError}
              disabled={loading}
              onChange={(e) => {
                const value = e.target.value;
                setSecret(value);
                setSecretError(validateSecretField(value));
              }}
              required={fieldsRequired}
              autoComplete="off"
            />
          </div>
          <div className="form__field">
            <FormInput
              label={t('settings.header')}
              name="header"
              type="text"
              value={header}
              error={headerError}
              disabled={loading}
              onChange={(e) => {
                const value = e.target.value;
                setHeader(value);
                setHeaderError(validateHeaderField(value));
              }}
              required={fieldsRequired}
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
                onChange={() => {
                  setDemo(!demo);
                }}
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
              disabled={saveDisabled}
              className="form__save-button"
              title={t('settings.save')}
            />
          </div>
        </form>
      </div>
    </div>
  );
});
