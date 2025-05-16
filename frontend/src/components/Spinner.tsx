import React from 'react';

import '@components/spinner.css';

interface SpinnerProps {
  size?: 'small' | 'medium' | 'large';
  className?: string;
  style?: React.CSSProperties;
  variant?: 'default' | 'blue' | 'colorful';
}

const Spinner: React.FC<SpinnerProps> = ({
  size = 'medium',
  className = '',
  style = {},
  variant = 'default',
}) => {
  const defaultStyle = {
    margin: '16px',
    ...style,
  };

  const variantClass = variant !== 'default' ? `spinner_${variant}` : '';

  return (
    <div
      className={`spinner spinner_${size} ${variantClass} ${className}`}
      style={defaultStyle}
    >
      <div className="spinner__circle">
        <div className="spinner__circle_gradient" />
        <div className="spinner__circle_inner" />
      </div>
    </div>
  );
};

Spinner.displayName = 'Spinner';

export default Spinner;
