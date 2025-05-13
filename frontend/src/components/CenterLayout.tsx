import React, { CSSProperties } from 'react';

import '@components/center.css';

export interface CenterLayoutProps {
  children: React.ReactNode;
  className?: string;
  style?: CSSProperties;
}

export const CenterLayout: React.FC<CenterLayoutProps> = ({ 
  children, 
  className = '', 
  style 
}) => {
  return (
    <div 
      className={`center-layout ${className}`} 
      style={style}
    >
      {children}
    </div>
  );
};
