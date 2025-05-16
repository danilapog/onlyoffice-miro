import React, { forwardRef } from 'react';

import '@components/label.css';

interface LabelProps extends React.LabelHTMLAttributes<HTMLLabelElement> {
  children: React.ReactNode;
  className?: string;
  htmlFor: string;
}

const Label = forwardRef<HTMLLabelElement, LabelProps>(
  ({ children, className = '', htmlFor, ...props }, ref) => {
    return (
      <label
        ref={ref}
        className={`generic-label ${className}`}
        htmlFor={htmlFor}
        {...props}
      >
        {children}
      </label>
    );
  }
);

Label.displayName = 'Label';

export default Label;
