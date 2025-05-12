import React, { forwardRef } from 'react';

import '@components/button.css';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  name: string;
  variant?: 'primary' | 'default';
};

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(({
  id,
  name,
  disabled,
  value,
  className = '',
  variant = 'default',
  onClick,
  ...props
}, ref) => {
  const realId = id || Math.random().toString(36).substring(2, 9);

  const handleClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    if (onClick) onClick(e);
  };

  return (
    <button
      id={realId}
      ref={ref}
      disabled={disabled}
      onClick={handleClick}
      className={`generic-button ${variant === 'primary' ? 'primary' : ''} ${className}`}
      {...props}
    >
      {name}
    </button>
  )
});
