import React, { forwardRef, useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { useManagerStore } from '@features/manager/stores/useManagerStore';
import { useFilesStore } from '@features/file/stores/useFileStore';
import { FormInput } from '@components/Input';
import { Button } from '@components/Button';
import { Label } from '@components/Label';
import { Select, SelectOption } from '@components/Select';

import '@features/manager/components/creator.css';

interface CreatorProps extends React.HTMLAttributes<HTMLDivElement> {
};

export const Creator = forwardRef<HTMLDivElement, CreatorProps>(({
  className,
  children,
  ...props
}, ref) => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const {
    selectedName,
    selectedType,
    loading,
    setSelectedName,
    setSelectedType,
    resetSelected,
    createFile,
    getSupportedTypes,
  } = useManagerStore();
  const {
    refreshDocuments
  } = useFilesStore();
  const [nameError, setNameError] = useState<string | undefined>(undefined);

  useEffect(() => {
    setSelectedType(getSupportedTypes()[0]);
  }, [getSupportedTypes, setSelectedType]);

  useEffect(() => {
    if (!selectedName || selectedName.trim() === '')
      setNameError(t('creation.errors.nameRequired'));
    else
      setNameError(undefined);
  }, [selectedName, t]);

  const fileTypeOptions: SelectOption[] = useMemo(() => {
    const supportedTypes = getSupportedTypes();
    return supportedTypes.map(type => ({
      value: type,
      label: t(`file.types.${type}`)
    }));
  }, [getSupportedTypes, t]);

  const formValid = selectedName && selectedName.trim() !== '' && selectedType && !nameError;

  const handleCreateFile = async () => {
    await createFile();
    resetSelected();
    await refreshDocuments();
  }

  return (
    <div
      ref={ref}
      className={`creator-container ${className}`}
      {...props}
    >
      <div className='creator-container__main'>
        <div className='creator-container_shifted'>
          <div className='creator-container__form'>
            <div className='creator-container__form-group'>
              <Label className='creator-container__label'>
                {t('creation.fileName')}<span className="form-input__label_required">*</span>
              </Label>
              <FormInput 
                className='creator-container__input'
                value={selectedName}
                onChange={(e) => setSelectedName(e.target.value)}
                disabled={loading}
                error={nameError}
              />
            </div>
            <div className='creator-container__form-group'>
              <Label className='creator-container__label'>{t('creation.fileType')}</Label>
              <Select 
                options={fileTypeOptions}
                value={selectedType}
                onChange={setSelectedType}
                className='creator-container__select'
                disabled={loading}
              />
            </div>
          </div>
        </div>
      </div>
      <div className='creator-container_shifted'>
        <Button 
          name={loading ? t('creation.creating') : t('creation.create')}
          variant='primary' 
          onClick={async () => {
            await handleCreateFile();
            navigate('/', { state: { isBack: true } });
          }}
          disabled={loading || !formValid}
        />
      </div>
    </div>
  );
});
