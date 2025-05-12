import React, { forwardRef } from 'react';

import '@components/label.css';

interface LabelProps extends React.LabelHTMLAttributes<HTMLLabelElement> {
  children: React.ReactNode;
  className?: string;
}

export const Label = forwardRef<HTMLLabelElement, LabelProps>(({
  children,
  className = '',
  ...props
}, ref) => {
  return (
    <label
      ref={ref}
      className={`generic-label ${className}`}
      {...props}
    >
      {children}
    </label>
  );
}); 