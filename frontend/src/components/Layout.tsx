import React, { ReactNode } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { Footer } from '@components/Footer';

import { useApplicationStore } from '@stores/useApplicationStore';

import '@components/layout.css';

interface LayoutProps {
  title: string;
  subtitle?: string;
  footerText: string;
  backTo?: string;
  settings?: boolean;
  help?: boolean;
  onSettings?: () => void;
  onHelp?: () => void;
  children: ReactNode;
}

export const Layout: React.FC<LayoutProps> = ({ 
    title,
    subtitle,
    footerText,
    backTo,
    settings,
    help,
    onSettings,
    onHelp = () => window.open('https://onlyoffice.com', '_blank'),
    children,
}) => {
  const { admin } = useApplicationStore();
  const navigate = useNavigate();
  
  if (!onSettings)
    onSettings = () => navigate("/settings");
  
  const showSettings = settings !== false && admin;
  return (
    <div className="layout-container">
      <div className="layout-header">
        <div className="layout-shifted">
          {backTo ? (
            <div className='layout-title'>
              <Link to={backTo} className='layout-title__back'>
                <img src="/arrow-left.svg" alt="Back" className="layout-title__back-icon" />
                <span className="layout-title__back-text" title={title}>{title}</span>
              </Link>
            </div>
          ) : (
            <div className="layout-title" title={title}>{title}</div>
          )}
          {subtitle && (
            <div className="layout-subtitle">
              {subtitle}
            </div>
          )}
        </div>
      </div>
      <div className="layout-main">
        {children}
      </div>
      <div className="layout-footer">
        <Footer text={footerText} settings={showSettings} help={help} onSettingsClick={onSettings} onHelpClick={onHelp} />
      </div>
    </div>
  );
}; 
