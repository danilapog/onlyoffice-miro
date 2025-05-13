import React, { forwardRef } from 'react';

import '@features/file/components/empty.css';

interface EmptyProps extends React.HTMLAttributes<HTMLDivElement> {
  title: string;
  subtitle: string;
};

export const Empty = forwardRef<HTMLDivElement, EmptyProps>(({
  title,
  subtitle,
  className,
  ...props
}, ref) => {
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
      <span className='empty-container__title'>{title}</span>
      <span className='empty-container__subtitle'>{subtitle}</span>
    </div>
  )
});
