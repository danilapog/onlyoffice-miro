import React, { forwardRef } from 'react';
import { useTranslation } from 'react-i18next';

import '@features/file/components/empty.css';

interface EmptyProps extends React.HTMLAttributes<HTMLDivElement> {
};

export const Empty = forwardRef<HTMLDivElement, EmptyProps>(({
  className,
  ...props
}, ref) => {
  const { t } = useTranslation();

  return (
    <div
      ref={ref}
      className='empty-container'
      {...props}
    >
      <img
        className='empty-container__icon'
        src='/nodocs.svg'
      />
      <span className='empty-container__title'>{t('empty.title')}</span>
      <span className='empty-container__subtitle'>{t('empty.subtitle')}</span>
    </div>
  )
});
